package f1

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirFirstTrue{})
}

type fnAirFirstTrue struct {
}

func (fnAirFirstTrue) Name() string {
	return "airfirsttrue"
}

func (fnAirFirstTrue) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeArray}, false
}

func (fnAirFirstTrue) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[1] {
		return false, nil
	}

	conditions, ok := params[1].([]interface{})
	if !ok {
		return false, errors.New("ERROR! Not an array!")
	}

	reading := params[0].(map[string]interface{})
	for index, element := range conditions {
		condition, ok := element.(map[string]interface{})
		if !ok {
			return false, errors.New("ERROR! Condition is not a map!")
		}
		result := true
		for key, value := range condition {
			if reading[key] != value {
				result = false
				break
			}
		}
		if result {
			return index, nil
		}
	}

	return -1, nil
}
