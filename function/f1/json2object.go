package f1

import (
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnJson2Object{})
}

type fnJson2Object struct {
}

func (fnJson2Object) Name() string {
	return "json2object"
}

func (fnJson2Object) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnJson2Object) Eval(params ...interface{}) (interface{}, error) {

	log.Debug("(fnJson2Object.Eval) params[0] : ", params[0])

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
		log.Warn("(fnJson2Object.Eval) err : ", err.Error())
		return nil, err
	}
	log.Debug("(fnJson2Object.Eval) rootObject : ", rootObject)

	return rootObject, nil
}
