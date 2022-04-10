package air

import (
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var airDataSelector = &fnAirDataSelector{}

func Test01(t *testing.T) {
	log.SetLogLevel(logger.DebugLevel)
	dataStr := "{\"gateway\":\"abc\",\"readings\":[{\"id\": \"51f57c16-fac5-4392-a3ae-b82fe3843e46\",\"origin\": 1644958988406173446,\"deviceName\": \"device1\",\"resourceName\": \"PipelineParameters\",\"profileName\": \"ValitaCell\",\"valueType\": \"Object\",\"objectValue\": {\"InputFileLocation\": \"/tmp/files/input/test-image1.tiff\",\"OutputFileFolder\": \"/tmp/files/output\",\"ModelParams\": {\"Brighten\": \"80\"},\"JobUpdateUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567\",\"PipelineStatusUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567\"}}], \"enriched\":[{\"producer\":\"PythonService1\",\"name\":\"Result\",\"value\":{\"id\": \"process:abc\", \"input1\": [[2, 1], [3, 4]], \"input2\": [[6, 5], [8, 7]], \"result\": [2, 1, 3, 4, 6, 5, 8, 7]}}]}"
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		panic(err)
	}

	gateway := data["gateway"]
	readints := data["readings"].([]interface{})[0]
	enriched := data["enriched"]

	result, _ := airDataSelector.Eval(gateway, readints, enriched, "@PythonService1..Result/input1[0][1]@")
	if 1 != int(result.(float64)) {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][1]") not matched!`)
	}

	result, _ = airDataSelector.Eval(gateway, readints, enriched, "@PythonService1..Result/result[4]@")
	if 6 != int(result.(float64)) {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/result[4]") not matched!`)
	}

	result, _ = airDataSelector.Eval(gateway, readints, enriched, "@PythonService1..Result/id@")
	if "process:abc" != result {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/id") not matched!`)
	}

	result, _ = airDataSelector.Eval(gateway, readints, enriched, "@f1..gateway@")
	if "abc" != result {
		t.Fatalf(`ExtractData(dataMap, "@f1..gateway@") not matched!`)
	}
	result, _ = airDataSelector.Eval(gateway, readints, enriched, "@f1..origin@")
	if 1644958988406173440 != int(result.(float64)) {
		t.Fatalf(`ExtractData(dataMap, "f1..readings/origin") not matched!`)
	}

	result, _ = airDataSelector.Eval(gateway, readints, enriched, "@f1..objectValue/InputFileLocation@")
	if "/tmp/files/input/test-image1.tiff" != result {
		t.Fatalf(`ExtractData(dataMap, "f1..objectValue/InputFileLocation") not matched!`)
	}

	result, _ = airDataSelector.Eval(gateway, readints, enriched, "@f1..objectValue/ModelParams/Brighten@")
	if "80" != result {
		t.Fatalf(`ExtractData(dataMap, "f1..objectValue/ModelParams/Brighten") not matched!`)
	}
}

func Test02(t *testing.T) {
	log.SetLogLevel(logger.DebugLevel)
	dataStr := "{\"gateway\":\"abc\",\"readings\":[{\"id\": \"51f57c16-fac5-4392-a3ae-b82fe3843e46\",\"origin\": 1644958988406173446,\"deviceName\": \"device1\",\"resourceName\": \"PipelineParameters\",\"profileName\": \"ValitaCell\",\"valueType\": \"Object\",\"objectValue\": {\"InputFileLocation\": \"/tmp/files/input/test-image1.tiff\",\"OutputFileFolder\": \"/tmp/files/output\",\"ModelParams\": {\"Brighten\": \"80\"},\"JobUpdateUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567\",\"PipelineStatusUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567\"}}], \"enriched\":[{\"producer\":\"PythonService1\",\"name\":\"Result\",\"value\":{\"id\": \"process:abc\", \"input1\": [[2, 1], [3, 4]], \"input2\": [[6, 5], [8, 7]], \"result\": [2, 1, 3, 4, 6, 5, 8, 7]}}]}"
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		panic(err)
	}

	gateway := data["gateway"]
	readints := data["readings"].([]interface{})[0]
	enriched := data["enriched"]

	result, _ := airDataSelector.Eval(gateway, readints, enriched, "{\"input_location\": \"@f1..objectValue/InputFileLocation@\",\"output_folder\": \"@f1..objectValue/OutputFileFolder@\", \"cellpose_batch_size\": 1}")
	if "{\"input_location\": \"/tmp/files/input/test-image1.tiff\",\"output_folder\": \"/tmp/files/output\", \"cellpose_batch_size\": 1}" != result {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][1]") not matched!`)
	}
}

func Test03(t *testing.T) {
	log.SetLogLevel(logger.DebugLevel)
	dataStr := "{\"gateway\":\"abc\",\"readings\":[{\"id\": \"51f57c16-fac5-4392-a3ae-b82fe3843e46\",\"origin\": 1644958988406173446,\"deviceName\": \"device1\",\"resourceName\": \"PipelineParameters\",\"profileName\": \"ValitaCell\",\"valueType\": \"Object\",\"objectValue\": {\"InputFileLocation\": \"/tmp/files/input/test-image1.tiff\",\"OutputFileFolder\": \"/tmp/files/output\",\"ModelParams\": {\"Brighten\": \"80\"},\"JobUpdateUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567\",\"PipelineStatusUrl\" : \"http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567\"}}]}"
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		panic(err)
	}

	gateway := data["gateway"]
	readints := data["readings"].([]interface{})[0]
	enriched := data["enriched"]

	result, _ := airDataSelector.Eval(gateway, readints, enriched, "{\"input_location\": \"@f1..objectValue/InputFileLocation@\",\"output_folder\": \"@f1..objectValue/OutputFileFolder@\", \"cellpose_batch_size\": 1}")
	if "{\"input_location\": \"/tmp/files/input/test-image1.tiff\",\"output_folder\": \"/tmp/files/output\", \"cellpose_batch_size\": 1}" != result {
		t.Fatalf(`ExtractDataAsString(dataMap, "PythonService1..Result/input1[0][1]") not matched!`)
	}
}
