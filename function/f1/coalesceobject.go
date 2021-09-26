package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnCoalesceObject{})
}

type fnCoalesceObject struct {
}

func (fnCoalesceObject) Name() string {
	return "coalesceobject"
}

func (fnCoalesceObject) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject}, true
}

func (fnCoalesceObject) Eval(params ...interface{}) (interface{}, error) {
	for _, param := range params {
		if nil != param {
			return param, nil
		}
	}

	return nil, nil
}
