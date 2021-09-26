package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnHasKey{})
}

type fnHasKey struct {
}

func (fnHasKey) Name() string {
	return "haskey"
}

func (fnHasKey) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeString}, false
}

func (fnHasKey) Eval(params ...interface{}) (interface{}, error) {
	object := params[0]
	key := params[1]
	if nil == object || nil == key {
		return false, nil
	}

	_, exist := (object.(map[string]interface{}))[key.(string)]

	return exist, nil
}
