package air

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFNObjectDeserializer_Eval(t *testing.T) {
	f := &fnObjectDeserializer{}
	objectstr := "�ggatewaycabchreadings��bidx$51f57c16-fac5-4392-a3ae-b82fe3843e46forigin�|	/jdeviceNamegdevice1lresourceNamerPipelineParameterskprofileNamejValitaCellivalueTypefObjectkobjectValue�lJobUpdateUrlxRhttp://Air-account_00001_valitacell_api_service:10108/api/v1/job/pipeline/1234/567qPipelineStatusUrlxThttp://Air-account_00001_valitacell_api_service:10108/api/v1/pipelineStatus/1234/567qInputFileLocationx!/tmp/files/input/test-image1.tiffpOutputFileFolderq/tmp/files/outputkModelParams�hBrightenb80"
	v, err := function.Eval(f, []byte(objectstr), nil)
	assert.Nil(t, err)
	fmt.Println("#### ", v)
}

func TestFNObjectDeserializer_Eval2(t *testing.T) {
	f := &fnObjectDeserializer{}
	objectstr := "{\"sample_type\":\"JSON\"}"
	v, err := function.Eval(f, []byte(objectstr), nil)
	assert.Nil(t, err)
	fmt.Println("#### ", v)
}
