package air

import (
	"encoding/json"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnObject2Json{})
}

type fnObject2Json struct {
}

func (fnObject2Json) Name() string {
	return "object2json"
}

func (fnObject2Json) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

func (fnObject2Json) Eval(params ...interface{}) (interface{}, error) {
	log.Info("(fnObject2Json.Eval) params[0] : ", params[0])
	object := params[0]
	if nil == object {
		return "{}", nil
	}
	jsonStr, err := json.Marshal(object)
	if nil != err {
		return "", err
	}
	return string(jsonStr), nil
}
