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

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

const (
	oOutputDataCollection = "OutputDataCollection"
)

type Settings struct {
	Targets string `md:"Targets"`
}

type Input struct {
	Consumer            string                 `md:"Consumer"`
	Producer            string                 `md:"Producer"`
	TargetData          map[string]interface{} `md:"TargetData"`
	InputDataCollection []interface{}          `md:"InputDataCollection"`
}

type Output struct {
	OutputDataCollection []interface{} `md:"OutputDataCollection"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Consumer":            i.Consumer,
		"Producer":            i.Producer,
		"TargetData":          i.TargetData,
		"InputDataCollection": i.InputDataCollection,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	ok := true
	i.Consumer, ok = values["Consumer"].(string)
	if !ok {
		return errors.New("Illegal Consumer type, expect string.")
	}
	i.Producer, ok = values["Producer"].(string)
	if !ok {
		return errors.New("Illegal Producer type, expect string.")
	}
	i.TargetData, ok = values["TargetData"].(map[string]interface{})
	if !ok {
		return errors.New("Illegal TargetData type, expect map[string]interface{}.")
	}
	i.InputDataCollection, ok = values["InputDataCollection"].([]interface{})
	if !ok {
		return errors.New("Illegal InputDataCollection type, expect []interface{}.")
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"OutputDataCollection": o.OutputDataCollection,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	ok := true
	o.OutputDataCollection, ok = values["OutputDataCollection"].([]interface{})
	if !ok {
		return errors.New("Illegal OutputDataCollection type, expect []interface{}.")
	}
	return nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

type Activity struct {
	dataTypes map[string]string
}

func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	dataTypes := make(map[string]string)
	var variablesDef []interface{}
	if err := json.Unmarshal([]byte(settings.Targets), &variablesDef); err != nil {
		return nil, err
	}
	for _, variableDef := range variablesDef {
		variableInfo := variableDef.(map[string]interface{})
		dataTypes[variableInfo["Name"].(string)] = variableInfo["Type"].(string)
	}

	activity := &Activity{
		dataTypes: dataTypes,
	}

	return activity, nil
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	log := ctx.Logger()
	log.Debug("(Eval) entering ........ ")
	defer log.Debug("(Eval) exit ........ ")

	input := &Input{}
	ctx.GetInputObject(input)

	producer := input.Producer
	log.Debug("[Eval]  producer : ", producer)

	consumer := input.Consumer
	log.Debug("[Eval]  consumer : ", consumer)

	targetData := input.TargetData
	log.Debug("[Eval]  targetData : ", targetData)

	inputDataCollection := input.InputDataCollection
	log.Debug("[Eval]  inputDataCollection : ", inputDataCollection)

	outputDataCollection := make([]interface{}, len(inputDataCollection))
	for index, data := range inputDataCollection {
		outputDataCollection[index] = data
	}

	var newValue interface{}
	for key, value := range targetData {
		log.Debug("[Eval]  key : ", key, ", value : ", value)
		dataType := "String"
		if nil != a.dataTypes && "" != a.dataTypes[key] {
			dataType = a.dataTypes[key]
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
	ctx.SetOutput(oOutputDataCollection, outputDataCollection)

	log.Debug("[Eval] exit ........ ")

	return true, nil
}
