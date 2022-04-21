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
	if 3 < len(keyElements) {
		return nil
	}
	subkeyElements := strings.Split(keyElements[2], "/")
	log.Debug("(ExtractData) real keyword : ", keyElements[2])
	var data interface{}
	data = dataMap[fmt.Sprintf("%s..%s", keyElements[0], subkeyElements[0])]
	log.Debug("(ExtractData) First level data : ", data)
	if 1 < len(subkeyElements) {
		if _, isMap := data.(map[string]interface{}); !isMap {
			if _, isArray := data.([]interface{}); !isArray {
				return nil
			}
		}
		subkey := fmt.Sprintf("root%s", strings.Replace(keyElements[2][len(subkeyElements[0]):], "/", ".", -1))
		log.Debug("(ExtractData) subkey : ", subkey)
		data = objectbuilder.LocateObject(data.(map[string]interface{}), subkey).(interface{})
	}
	log.Debug("(ExtractData) data : ", data)
	return data
}

func ExtractDataAsString(dataMap map[string]interface{}, keyword string) string {
	data := ExtractData(dataMap, keyword)
	dataType := reflect.ValueOf(data).Kind()
	log.Debug("(ExtractDataAsString) dataType : ", dataType.String())
	if reflect.String == dataType {
		return strings.ReplaceAll(data.(string), "\"", "\\\"")
	} else if reflect.Map == dataType || reflect.Array == dataType || reflect.Slice == dataType {
		jsonBuf, _ := json.Marshal(data)
		log.Debug("(ExtractDataAsString) string(jsonBuf) : ", string(jsonBuf))
		return fmt.Sprintf("%v", string(jsonBuf))
	}
	return fmt.Sprintf("%v", data)
}
