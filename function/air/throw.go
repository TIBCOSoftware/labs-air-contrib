package air

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnThrow{})
}

type fnThrow struct {
}

func (fnThrow) Name() string {
	return "throw"
}

func (fnThrow) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBool, data.TypeString}, false
}

func (fnThrow) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] {
		return false, nil
	}

	if params[0].(bool) {
		return true, errors.New(params[1].(string))
	}
	return false, nil
}
