/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

/*
	{
		"imports": [],
		"name": "ProjectAirApplication",
		"description": "",
		"version": "1.0.0",
		"type": "flogo:app",
		"appModel": "1.1.1",
		"feVersion": "2.9.0",
		"triggers": [],
		"resources": [],
		"properties": [],
		"connections": {},
		"contrib": "",
		"fe_metadata": ""
	}
*/

package dataselector

import (
	"errors"
	"fmt"
	"sync"

	kwr "github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/keywordreplace"
	"github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/util"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("labs-lc-activity-DataSelector")

var initialized bool = false

const (
	sLeftToken      = "leftToken"
	sRightToken     = "rightToken"
	sVariablesDef   = "variablesDef"
	sTargets        = "targets"
	iVariable       = "Variables"
	iDataCollection = "DataCollection"
	oExtractedData  = "ExtractedData"
)

type DataSelector struct {
	metadata  *activity.Metadata
	mux       sync.Mutex
	selectors map[string]map[string]string
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aDataSelector := &DataSelector{
		metadata:  metadata,
		selectors: make(map[string]map[string]string),
	}

	return aDataSelector
}

func (a *DataSelector) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *DataSelector) Eval(context activity.Context) (done bool, err error) {

	log.Info("[DataSelector:Eval] entering ........ ")
	defer log.Info("[DataSelector:Eval] exit ........ ")

	inputDataCollection, ok := context.GetInput(iDataCollection).([]interface{})
	if !ok {
		return false, errors.New("Invalid inputDataCollection ... ")
	}

	selector, err := a.getSelector(context)
	log.Debug("[DataSelector:Eval] selector : ", selector)
	if nil != err {
		return true, err
	}

	mapping := make(map[string]interface{})
	for _, data := range inputDataCollection {
		producer := data.(map[string]interface{})["producer"]
		if nil == producer {
			producer = ""
		}
		consumer := data.(map[string]interface{})["consumer"]
		if nil == consumer {
			consumer = ""
		}
		name := data.(map[string]interface{})["name"]
		if nil == name {
			name = ""
		}

		log.Debug("[DataSelector:Eval] data key : ", fmt.Sprintf("%s.%s.%s", producer, consumer, name))

		//////
		// TO DO : Ignore consumer for now
		/////
		//mapping[fmt.Sprintf("%s.%s.%s", producer, consumer, name)] = data
		mapping[fmt.Sprintf("%s..%s", producer, name)] = data
	}

	log.Debug("[DataSelector:Eval] mapping : ", mapping)

	pathMapper := a.getVariableMapper(context)
	variable := context.GetInput(iVariable)
	log.Debug("[DataSelector:Eval] pathMapper : ", pathMapper)
	log.Debug("[DataSelector:Eval] variable : ", variable)
	extractedData := make(map[string]interface{})
	for key, name := range selector {
		if nil != variable && nil != pathMapper {
			key = pathMapper.Replace(key, variable.(map[string]interface{}))
			log.Info("[DataSelector:Eval] key : ", key)
		}
		if nil != mapping[key] {
			extractedData[name] = mapping[key].(map[string]interface{})["value"]
			log.Debug("[DataSelector:Eval] value found for key = ", key, ", value = ", mapping[key])
		} else {
			log.Warn("[DataSelector:Eval] value not found for key = ", key)
		}
	}

	log.Debug("[DataSelector:Eval]  oExtractedData : ", extractedData)
	context.SetOutput(oExtractedData, extractedData)

	return true, nil
}

func (a *DataSelector) getSelector(ctx activity.Context) (map[string]string, error) {
	myId := util.ActivityId(ctx)
	selector := a.selectors[myId]
	if nil == selector {
		a.mux.Lock()
		defer a.mux.Unlock()
		selector = a.selectors[myId]
		if nil == selector {

			variables := make(map[string]interface{})
			variablesDef, _ := ctx.GetSetting(sVariablesDef)
			log.Debug("Processing handlers : variablesDef = ", variablesDef)
			for _, variableDef := range variablesDef.([]interface{}) {
				variableInfo := variableDef.(map[string]interface{})
				variables[variableInfo["Name"].(string)] = variableInfo["Type"].(string)
			}

			selector = make(map[string]string)
			targetsDef, ok := ctx.GetSetting(sTargets)
			log.Debug("[DataSelector:getSelector] Processing handlers : sTargets = ", targetsDef)
			if ok && nil != targetsDef {
				for _, targetDef := range targetsDef.([]interface{}) {
					targetInfo := targetDef.(map[string]interface{})
					log.Debug("[DataSelector:getSelector] targetInfo = ", targetInfo)
					filedMatch := targetInfo["FieldMatch"].(string)
					selector[filedMatch] = targetInfo["Name"].(string)
					log.Debug("[DataSelector:getSelector] selector = ", selector)
				}
			}

			a.selectors[myId] = selector
		}
		log.Debug("[DataSelector:getSelector] selector = ", selector)
	}
	return selector, nil
}

func (a *DataSelector) getVariableMapper(ctx activity.Context) *kwr.KeywordMapper {
	lefttoken, exist := ctx.GetSetting(sLeftToken)
	if !exist {
		return nil
	}
	righttoken, exist := ctx.GetSetting(sRightToken)
	if !exist {
		return nil
	}
	return kwr.NewKeywordMapper("", lefttoken.(string), righttoken.(string))
}
