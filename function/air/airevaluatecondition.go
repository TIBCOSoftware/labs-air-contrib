package air

import (
	"errors"
	"fmt"

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
	return []data.Type{data.TypeString, data.TypeObject, data.TypeArray, data.TypeArray, data.TypeInt}, false
}

func (fnAirEvaluateCondition) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[3] {
		size := params[4].(int)
		results := make([]bool, size)
		for index := 0; index < size; index++ {
			results[index] = false
		}
		return results, nil
	}

	conditions, ok := params[3].([]interface{})
	if !ok {
		return false, errors.New("ERROR! Not an array!")
	}

	data := make(map[string]interface{})
	data["gateway"] = params[0]
	if nil != params[1] {
		for key, value := range params[1].(map[string]interface{}) {
			data[key] = value
		}
	}
	if nil != params[2] {
		for _, element := range params[2].([]interface{}) {
			enrichedElement := element.(map[string]interface{})
			data[fmt.Sprintf("%s.%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
		}
	}

	results := make([]bool, len(conditions))
	for index, element := range conditions {
		condition, ok := element.(map[string]interface{})
		if !ok {
			return false, errors.New("ERROR! Condition is not a map!")
		}
		results[index] = true
		for key, value := range condition {
			if data[key] != value {
				results[index] = false
				break
			}
		}
	}

	return results, nil
}
