package air

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirEvaluateCondition{})
}

type fnAirEvaluateCondition struct {
}

func (fnAirEvaluateCondition) Name() string {
	return "airevaluatecondition"
}

func (fnAirEvaluateCondition) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeArray}, false
}

func (fnAirEvaluateCondition) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[1] {
		return make([]bool, 0), nil
	}

	conditions, ok := params[1].([]interface{})
	if !ok {
		return false, errors.New("ERROR! Not an array!")
	}

	reading := params[0].(map[string]interface{})
	results := make([]bool, len(conditions))
	for index, element := range conditions {
		condition, ok := element.(map[string]interface{})
		if !ok {
			return false, errors.New("ERROR! Condition is not a map!")
		}
		results[index] = true
		for key, value := range condition {
			if reading[key] != value {
				results[index] = false
				break
			}
		}
	}

	return results, nil
}
