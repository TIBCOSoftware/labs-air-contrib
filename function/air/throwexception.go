package air

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnThrowException{})
}

type fnThrowException struct {
}

func (fnThrowException) Name() string {
	return "throwexception"
}

func (fnThrowException) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBool, data.TypeString}, false
}

func (fnThrowException) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] {
		return false, nil
	}

	if params[0].(bool) {
		return true, errors.New(params[1].(string))
	}
	return false, nil
}
