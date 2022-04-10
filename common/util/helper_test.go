package util

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func TestExtractData(t *testing.T) {
	log.SetLogLevel(logger.DebugLevel)
	dataStr := "{\"gateway\":\"abc\",\"readings\":[{\"id\": \"51f57c16-fac5-4392-a3ae-b82fe3843e46\",\"origin\": 1644958988406173446,\"deviceName\": \"device1\",\"resourceName\": \"PipelineParameters\",\"profileName\": \"ValitaCell\",\"valueType\": \"Object\",\"objectValue\": {\"InputFileLocation\": \"/tmp/files/input/test-image1.tiff\",\"OutputFileFolder\": \"/tmp/files/output\",\"ModelParams\": {\"Brighten\": \"80\"},\"JobUpdateUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567\",\"PipelineStatusUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567\"}}], \"enriched\":[{\"producer\":\"PythonService1\",\"name\":\"Result\",\"value\":{\"id\": \"process:abc\", \"input1\": [[2, 1], [3, 4]], \"input2\": [[6, 5], [8, 7]], \"result\": [2, 1, 3, 4, 6, 5, 8, 7]}}]}"
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		panic(err)
	}

	dataMap := make(map[string]interface{})
	dataMap["f1..gateway"] = data["gateway"]                                               // gateway
	for key, value := range data["readings"].([]interface{})[0].(map[string]interface{}) { // reading
		dataMap[fmt.Sprintf("f1..%s", key)] = value
	}
	if nil != data["enriched"] { // enriched
		for _, element := range data["enriched"].([]interface{}) {
			enrichedElement := element.(map[string]interface{})
			dataMap[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
		}
	}

	if "abc" != ExtractData(dataMap, "f1..gateway") {
		t.Fatalf(`ExtractData(dataMap, "f1..gateway") not matched!`)
	}

	if 1644958988406173440 != int(ExtractData(dataMap, "f1..origin").(float64)) {
		t.Fatalf(`ExtractData(dataMap, "f1..readings/origin") not matched!`)
	}

	if "/tmp/files/input/test-image1.tiff" != ExtractData(dataMap, "f1..objectValue/InputFileLocation") {
		t.Fatalf(`ExtractData(dataMap, "f1..objectValue/InputFileLocation") not matched!`)
	}

	if "80" != ExtractData(dataMap, "f1..objectValue/ModelParams/Brighten") {
		t.Fatalf(`ExtractData(dataMap, "f1..objectValue/ModelParams/Brighten") not matched!`)
	}
}

func TestExtractDataAsString(t *testing.T) {
	log.SetLogLevel(logger.DebugLevel)
	dataStr := "{\"gateway\":\"abc\",\"readings\":[{\"id\": \"51f57c16-fac5-4392-a3ae-b82fe3843e46\",\"origin\": 1644958988406173446,\"deviceName\": \"device1\",\"resourceName\": \"PipelineParameters\",\"profileName\": \"ValitaCell\",\"valueType\": \"Object\",\"objectValue\": {\"InputFileLocation\": \"/tmp/files/input/test-image1.tiff\",\"OutputFileFolder\": \"/tmp/files/output\",\"ModelParams\": {\"Brighten\": \"80\"},\"JobUpdateUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567\",\"PipelineStatusUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567\"}}], \"enriched\":[{\"producer\":\"PythonService1\",\"name\":\"Result\",\"value\":{\"id\": \"process:abc\", \"input1\": [[2, 1], [3, 4]], \"input2\": [[6, 5], [8, 7]], \"result\": [2, 1, 3, 4, 6, 5, 8, 7]}}]}"
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		panic(err)
	}

	dataMap := make(map[string]interface{})
	dataMap["f1..gateway"] = data["gateway"]                                               // gateway
	for key, value := range data["readings"].([]interface{})[0].(map[string]interface{}) { // reading
		dataMap[fmt.Sprintf("f1..%s", key)] = value
	}
	if nil != data["enriched"] { // enriched
		for _, element := range data["enriched"].([]interface{}) {
			enrichedElement := element.(map[string]interface{})
			dataMap[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
		}
	}

	if "{\"id\":\"process:abc\",\"input1\":[[2,1],[3,4]],\"input2\":[[6,5],[8,7]],\"result\":[2,1,3,4,6,5,8,7]}" != ExtractDataAsString(dataMap, "PythonService1..Result") {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result") not matched!`)
	}

	if "[[2,1],[3,4]]" != ExtractDataAsString(dataMap, "PythonService1..Result/input1[]") {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/input1[]") not matched!`)
	}

	if "[2,1]" != ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][]") {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][]") not matched!`)
	}

	if "1" != ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][1]") {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][1]") not matched!`)
	}

	if "process:abc" != ExtractDataAsString(dataMap, "PythonService1..Result/id") {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/id") not matched!`)
	}

	if "abc" != ExtractDataAsString(dataMap, "f1..gateway") {
		t.Fatalf(`ExtractDataAsString(dataMap, "f1..gateway") not matched!`)
	}

	log.Info(ExtractDataAsString(dataMap, "f1..origin"))
	if "1.6449589884061734e+18" != ExtractDataAsString(dataMap, "f1..origin") {
		t.Fatalf(`ExtractDataAsString(dataMap, "f1..readings/origin") not matched!`)
	}

	if "/tmp/files/input/test-image1.tiff" != ExtractDataAsString(dataMap, "f1..objectValue/InputFileLocation") {
		t.Fatalf(`ExtractDataAsString(dataMap, "f1..objectValue/InputFileLocation") not matched!`)
	}

	if "80" != ExtractDataAsString(dataMap, "f1..objectValue/ModelParams/Brighten") {
		t.Fatalf(`ExtractDataAsString(dataMap, "f1..objectValue/ModelParams/Brighten") not matched!`)
	}
}
