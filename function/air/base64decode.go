package air

import (
	"encoding/base64"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnBase64Decode{})
}

type fnBase64Decode struct {
}

func (fnBase64Decode) Name() string {
	return "base64decode"
}

func (fnBase64Decode) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBool, data.TypeString}, false
}

func (fnBase64Decode) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] || !params[0].(bool) {
		return []byte(params[1].(string)), nil
	}

	value, err := base64.StdEncoding.DecodeString(params[1].(string))
	if nil != err {
		return nil, err
	}
	return value, nil
}
