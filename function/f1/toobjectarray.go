package f1

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

type ToObjectArray struct {
}

func init() {
	function.Register(&ToObjectArray{})
}

func (s *ToObjectArray) Name() string {
	return "toobjectarray"
}

func (s *ToObjectArray) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

func (s *ToObjectArray) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] {
		return nil, nil
	}
	out, ok := params[0].([]interface{})
	if !ok {
		return nil, errors.New("Unable to cast to array!")
	} else {
		return out, nil
	}
}
