package air

import (
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnEnsureJson{})
}

type fnEnsureJson struct {
}

func (fnEnsureJson) Name() string {
	return "ensurejson"
}

func (fnEnsureJson) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnEnsureJson) Eval(params ...interface{}) (interface{}, error) {
	log.Debug("(fnEnsureJson.Eval) params[0] : ", params[0])

	in, ok := params[0].(string)
	if !ok {
		err := fmt.Errorf("Illegal parameter!")
		log.Error("(fnEnsureJson.Eval) err : ", err.Error())
		return in, err
	}

	var rootObject interface{}
	err := json.Unmarshal([]byte(in), &rootObject)
	if nil != err {
		log.Error("(fnEnsureJson.Eval) build object, err : ", err.Error())
		return in, err
	}

	log.Debug("(fnEnsureJson.Eval) rootObject : ", rootObject)

	jsonStr, err := json.Marshal(rootObject)
	if nil != err {
		log.Error("(fnEnsureJson.Eval) object to string, err : ", err.Error())
		return in, err
	}
	return string(jsonStr), nil
}
