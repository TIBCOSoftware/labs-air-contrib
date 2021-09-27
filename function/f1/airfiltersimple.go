package f1

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirFilterSimple{})
}

type fnAirFilterSimple struct {
}

func (fnAirFilterSimple) Name() string {
	return "airfiltersimple"
}

func (fnAirFilterSimple) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeArray}, false
}

func (fnAirFilterSimple) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[1] {
		return false, nil
	}

	conditions, ok := params[1].([]interface{})
	if !ok {
		return false, errors.New("ERROR! Not an array!")
	}

	reading := params[0].(map[string]interface{})
	for _, element := range conditions {
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
			return true, nil
		}
	}

	return false, nil
}
