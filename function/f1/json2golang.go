package f1

import (
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnJson2Golang{})
}

type fnJson2Golang struct {
}

func (fnJson2Golang) Name() string {
	return "json2golang"
}

func (fnJson2Golang) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnJson2Golang) Eval(params ...interface{}) (interface{}, error) {
	in, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("Illegal parameter!")
	}
	if "" == in {
		return make(map[string]interface{}), nil
	}
	var rootObject interface{}
	err := json.Unmarshal([]byte(in), &rootObject)
	if nil != err {
		return nil, err
	}

	return rootObject, nil
}
