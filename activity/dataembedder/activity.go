/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package dataembedder

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-lightcrane-contrib/common/util"
)

var log = logger.GetLogger("tibco-f1-DataEmbedder")

var initialized bool = false

const (
	sTargets              = "Targets"
	iConsumer             = "Consumer"
	iProducer             = "Producer"
	iTargetData           = "TargetData"
	iInputDataCollection  = "InputDataCollection"
	oOutputDataCollection = "OutputDataCollection"
)

type DataEmbedder struct {
	metadata  *activity.Metadata
	mux       sync.Mutex
	dataTypes map[string]map[string]string
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aDataEmbedder := &DataEmbedder{
		metadata:  metadata,
		dataTypes: make(map[string]map[string]string),
	}

	return aDataEmbedder
}

func (a *DataEmbedder) Metadata() *activity.Metadata {

	return a.metadata
}

func (a *DataEmbedder) Eval(context activity.Context) (done bool, err error) {

	log.Debug("[Eval] entering ........ ")

	producer, ok := context.GetInput(iProducer).(string)
	if !ok {
		return false, errors.New("Invalid producer ... ")
	}
	log.Debug("[Eval]  producer : ", producer)

	consumer, ok := context.GetInput(iConsumer).(string)
	if !ok {
		return false, errors.New("Invalid consumer ... ")
	}
	log.Debug("[Eval]  consumer : ", consumer)

	targetData, ok := context.GetInput(iTargetData).(map[string]interface{})
	if !ok {
		return false, errors.New("Invalid targetData ... ")
	}
	log.Debug("[Eval]  targetData : ", targetData)

	inputDataCollection, ok := context.GetInput(iInputDataCollection).([]interface{})
	if !ok {
		return false, errors.New("Invalid inputDataCollection ... ")
	}
	log.Debug("[Eval]  inputDataCollection : ", inputDataCollection)

	outputDataCollection := make([]interface{}, len(inputDataCollection))
	for index, data := range inputDataCollection {
		outputDataCollection[index] = data
	}

	dataTypes, _ := a.getEnrichedDataType(context)
	var newValue interface{}
	for key, value := range targetData {
		log.Debug("[Eval]  key : ", key, ", value : ", value)
		dataType := "String"
		if nil != dataTypes && "" != dataTypes[key] {
			dataType = dataTypes[key]
		}

		log.Debug("[Eval]  dataType 01 : ", dataType)
		if "String" == dataType {
			var objectValue map[string]interface{}
			err := json.Unmarshal([]byte(value.(string)), &objectValue)
			log.Debug("[Eval]  objectValue : ", objectValue)
			if nil != err {
				log.Debug("[Eval]  Not object type : ", err.Error())
				newValue = value
			} else {
				newValue = objectValue
				dataType = "Object"
			}
		} else {
			newValue = value
		}

		log.Debug("[Eval]  dataType 02 : ", dataType)
		log.Debug("[Eval]  newValue : ", newValue)
		log.Debug("[Eval]  golang dataType : ", reflect.ValueOf(newValue).Kind().String())
		if nil != value {
			outputDataCollection = append(outputDataCollection, map[string]interface{}{
				"producer": producer,
				"consumer": consumer,
				"name":     key,
				"value":    newValue,
				"type":     dataType,
			})
		}
	}

	log.Debug("[Eval]  oOutputDataCollection : ", outputDataCollection)
	context.SetOutput(oOutputDataCollection, outputDataCollection)

	log.Debug("[Eval] exit ........ ")

	return true, nil
}

func (a *DataEmbedder) getEnrichedDataType(ctx activity.Context) (map[string]string, error) {
	myId := util.ActivityId(ctx)
	dataTypes := a.dataTypes[myId]

	if nil == dataTypes {
		a.mux.Lock()
		defer a.mux.Unlock()
		if nil == dataTypes {
			dataTypes = make(map[string]string)
			variablesDef, _ := ctx.GetSetting(sTargets)
			log.Debug("[getEnrichedDataType] enriched data def = ", variablesDef)
			for _, variableDef := range variablesDef.([]interface{}) {
				variableInfo := variableDef.(map[string]interface{})
				dataTypes[variableInfo["Name"].(string)] = variableInfo["Type"].(string)
			}
			a.dataTypes[myId] = dataTypes
		}
	}
	return dataTypes, nil
}
