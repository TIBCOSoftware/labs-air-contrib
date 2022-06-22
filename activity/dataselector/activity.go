/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package dataselector

import (
	"encoding/json"
	"errors"
	"fmt"

	kwr "github.com/TIBCOSoftware/labs-lightcrane-contrib/common/keywordreplace"
	//	"github.com/TIBCOSoftware/labs-lightcrane-contrib/common/util"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

const (
	oExtractedData = "ExtractedData"
)

type Settings struct {
	LeftToken    string `md:"leftToken"`
	RightToken   string `md:"rightToken"`
	VariablesDef string `md:"variablesDef"`
	Targets      string `md:"targets"`
}

type Input struct {
	Variables      map[string]interface{} `md:"Variables"`
	DataCollection []interface{}          `md:"DataCollection"`
}

type Output struct {
	ExtractedData interface{} `md:"ExtractedData"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Variables":      i.Variables,
		"DataCollection": i.DataCollection,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	ok := true
	i.Variables, ok = values["Variables"].(map[string]interface{})
	if !ok {
		return errors.New("Illegal Variables type, expect map[string]interface{}.")
	}
	i.DataCollection, ok = values["DataCollection"].([]interface{})
	if !ok {
		return errors.New("Illegal DataCollection type, expect []interface{}.")
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ExtractedData": o.ExtractedData,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.ExtractedData = values["ExtractedData"].(interface{})
	return nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

type Activity struct {
	selector      map[string]string
	keywordMapper *kwr.KeywordMapper
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	var variablesDef []interface{}
	if err := json.Unmarshal([]byte(settings.VariablesDef), &variablesDef); err != nil {
		return nil, err
	}

	variables := make(map[string]interface{})
	for _, variableDef := range variablesDef {
		variableInfo := variableDef.(map[string]interface{})
		variables[variableInfo["Name"].(string)] = variableInfo["Type"].(string)
	}

	var targetsDef []interface{}
	if err := json.Unmarshal([]byte(settings.Targets), &targetsDef); err != nil {
		return nil, err
	}

	selector := make(map[string]string)
	if nil != targetsDef {
		for _, targetDef := range targetsDef {
			targetInfo := targetDef.(map[string]interface{})
			filedMatch := targetInfo["FieldMatch"].(string)
			selector[filedMatch] = targetInfo["Name"].(string)
		}
	}

	var keywordMapper *kwr.KeywordMapper
	lefttoken := settings.LeftToken
	righttoken := settings.RightToken
	if "" != lefttoken && "" != righttoken {
		keywordMapper = kwr.NewKeywordMapper("", lefttoken, righttoken)
	}

	activity := &Activity{
		selector:      selector,
		keywordMapper: keywordMapper,
	}

	return activity, nil
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	log := ctx.Logger()
	log.Info("[DataSelector:Eval] entering ........ ")
	defer log.Info("[DataSelector:Eval] exit ........ ")

	input := &Input{}
	ctx.GetInputObject(input)

	log.Debug("[DataSelector:Eval] selector : ", a.selector)

	mapping := make(map[string]interface{})
	for _, data := range input.DataCollection {
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

	log.Debug("[DataSelector:Eval] pathMapper : ", a.keywordMapper)
	log.Debug("[DataSelector:Eval] variable : ", input.Variables)
	extractedData := make(map[string]interface{})
	for key, name := range a.selector {
		if nil != input.Variables && nil != a.keywordMapper {
			key = a.keywordMapper.Replace(key, input.Variables)
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
	ctx.SetOutput(oExtractedData, extractedData)

	return true, nil
}
