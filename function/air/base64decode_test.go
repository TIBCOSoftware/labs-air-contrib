package air

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFnBase64DecodeTrue(t *testing.T) {
	f := &fnBase64Decode{}
	v, err := function.Eval(f, true, "bXF0dF9hZG1pbg==")
	assert.Nil(t, err)
	assert.Equal(t, []byte("mqtt_admin"), v)
	fmt.Println("#### ", v)
}

func TestFnBase64DecodeFalse(t *testing.T) {
	f := &fnBase64Decode{}
	v, err := function.Eval(f, false, "bXF0dF9hZG1pbg==")
	assert.Nil(t, err)
	assert.Equal(t, []byte("bXF0dF9hZG1pbg=="), v)
	fmt.Println("#### ", v)
}
