package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnGetSubobject{})
}

type fnGetSubobject struct {
}

func (fnGetSubobject) Name() string {
	return "getsubobject"
}

func (fnGetSubobject) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeString}, false
}

func (fnGetSubobject) Eval(params ...interface{}) (interface{}, error) {
	log.Info("(fnGetSubobject.Eval) params[0] : ", params[0], ", params[1] : ", params[1])
	if nil == params[0] || nil == params[1] {
		return nil, nil
	}
	object := params[0].(map[string]interface{})
	key := params[1].(string)
	result := object[key]
	log.Info("(fnGetSubobject.Eval) object : ", object, ", key : ", key, ", ", ", result : ", result)

	return result, nil
}
