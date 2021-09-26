package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnIsEmpty{})
}

type fnIsEmpty struct {
}

func (fnIsEmpty) Name() string {
	return "isempty"
}

func (fnIsEmpty) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject}, false
}

func (fnIsEmpty) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] {
		return true, nil
	} else {
		return 0 == len(params[0].(map[string]interface{})), nil
	}
}
