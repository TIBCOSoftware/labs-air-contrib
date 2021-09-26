package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

type Ternary struct {
}

func init() {
	function.Register(&Ternary{})
}

func (s *Ternary) Name() string {
	return "ternary"
}

func (s *Ternary) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBool, data.TypeAny, data.TypeAny}, false
}

func (s *Ternary) Eval(params ...interface{}) (interface{}, error) {
	if params[0].(bool) {
		log.Debug("(fnTernary) True, return = ", params[1])
		return params[1], nil
	} else {
		log.Debug("(fnTernary) False, return = ", params[2])
		return params[2], nil
	}
}
