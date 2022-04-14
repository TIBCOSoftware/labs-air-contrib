package air

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFNObjectSerializer_Eval1(t *testing.T) {
	f := &fnObjectSerializer{}
	object := map[string]interface{}{
		"gateway": "abc",
		"readings": []interface{}{
			map[string]interface{}{
				"id":           "51f57c16-fac5-4392-a3ae-b82fe3843e46",
				"origin":       1644958988406173446,
				"deviceName":   "device1",
				"resourceName": "PipelineParameters",
				"profileName":  "ValitaCell",
				"valueType":    "Object",
				"objectValue": map[string]interface{}{
					"InputFileLocation": "/tmp/files/input/test-image1.tiff",
					"OutputFileFolder":  "/tmp/files/output",
					"ModelParams": map[string]interface{}{
						"Brighten": "80",
					},
					"JobUpdateUrl":      "http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567",
					"PipelineStatusUrl": "http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567",
				},
			},
		},
	}

	v, err := function.Eval(f, object, "JSON")
	assert.Nil(t, err)
	fmt.Println("#### ", string(v.([]byte)))
}

func TestFNObjectSerializer_Eval2(t *testing.T) {
	f := &fnObjectSerializer{}
	object := map[string]interface{}{
		"gateway": "abc",
		"readings": []interface{}{
			map[string]interface{}{
				"id":           "51f57c16-fac5-4392-a3ae-b82fe3843e46",
				"origin":       1644958988406173446,
				"deviceName":   "device1",
				"resourceName": "PipelineParameters",
				"profileName":  "ValitaCell",
				"valueType":    "Object",
				"objectValue": map[string]interface{}{
					"InputFileLocation": "/tmp/files/input/test-image1.tiff",
					"OutputFileFolder":  "/tmp/files/output",
					"ModelParams": map[string]interface{}{
						"Brighten": "80",
					},
					"JobUpdateUrl":      "http://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567",
					"PipelineStatusUrl": "http://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567",
				},
			},
		},
	}

	v, err := function.Eval(f, object, "CBOR")
	assert.Nil(t, err)
	fmt.Println("#### ", string(v.([]byte)))

	f2 := &fnObjectDeserializer{}
	v2, err := function.Eval(f2, v, nil)
	assert.Nil(t, err)
	fmt.Println("#### ", v2)

}
