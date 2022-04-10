/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-lightcrane-contrib/common/objectbuilder"
)

var log = logger.GetLogger("tibco-air_helper")

func ExtractData(dataMap map[string]interface{}, keyword string) interface{} {
	keyElements := strings.Split(keyword, ".")
	subkeyElements := strings.Split(keyElements[2], "/")
	log.Debug("(ExtractData) real keyword : ", keyElements[2])
	var data interface{}
	data = dataMap[fmt.Sprintf("%s..%s", keyElements[0], subkeyElements[0])]
	log.Debug("(ExtractData) First level data : ", data)
	if len(subkeyElements) > 1 {
		if _, isMap := data.(map[string]interface{}); !isMap {
			return nil
		}
		subkey := fmt.Sprintf("root%s", strings.Replace(keyElements[2][len(subkeyElements[0]):], "/", ".", -1))
		log.Debug("(ExtractData) subkey : ", subkey)
		data = objectbuilder.LocateObject(data.(map[string]interface{}), subkey).(interface{})
	}
	log.Debug("(ExtractDataAsString) data : ", data)
	return data
}

func ExtractDataAsString(dataMap map[string]interface{}, keyword string) string {
	keyElements := strings.Split(keyword, ".")
	subkeyElements := strings.Split(keyElements[2], "/")
	log.Debug("(ExtractDataAsString) real keyword : ", keyElements[2])
	var data interface{}
	data = dataMap[fmt.Sprintf("%s..%s", keyElements[0], subkeyElements[0])]
	log.Debug("(ExtractDataAsString) real data : ", data)

	dataType := reflect.ValueOf(data).Kind()
	log.Debug("(ExtractDataAsString) dataType : ", dataType.String())
	if reflect.String == dataType {
		return strings.ReplaceAll(data.(string), "\"", "\\\"")
	} else if reflect.Map == dataType {
		if len(subkeyElements) > 1 {
			log.Debug("(ExtractDataAsString) keyElements[2] : ", keyElements[2])
			subkey := fmt.Sprintf("root%s", strings.Replace(keyElements[2][len(subkeyElements[0]):], "/", ".", -1))
			log.Debug("(ExtractDataAsString) subkey : ", subkey)
			data = objectbuilder.LocateObject(data.(map[string]interface{}), subkey).(interface{})
			log.Debug("(ExtractDataAsString) data : ", data)
		}
		realDataType := reflect.ValueOf(data).Kind()
		log.Debug("(ExtractDataAsString) realDataType : ", realDataType.String())
		if reflect.Map == realDataType || reflect.Array == realDataType || reflect.Slice == realDataType {
			jsonBuf, _ := json.Marshal(data)
			log.Debug("(ExtractDataAsString) string(jsonBuf) : ", string(jsonBuf))
			return fmt.Sprintf("%v", string(jsonBuf))
		} else {
			log.Debug("(ExtractDataAsString) data.(string) : ", data.(string))
			return strings.ReplaceAll(data.(string), "\"", "\\\"")
		}
	} else if reflect.Array == dataType {
		jsonBuf, _ := json.Marshal(data)
		return fmt.Sprintf("%v", string(jsonBuf))
	}
	return fmt.Sprintf("%v", data)
}
