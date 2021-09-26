package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnModifyObject{})
}

type fnModifyObject struct {
}

func (fnModifyObject) Name() string {
	return "modifyobject"
}

func (fnModifyObject) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeString, data.TypeAny}, false
}

func (fnModifyObject) Eval(params ...interface{}) (interface{}, error) {
	log.Debug("(fnModifyObject) params[0] object = ", params[0])
	log.Info("(fnModifyObject) params[1] key = ", params[1])
	log.Debug("(fnModifyObject) params[2] value = ", params[2])

	if nil == params[0] || nil == params[1] {
		log.Warn("(fnModifyObject) params[0] >>>>>>>>>>", params[0])
		log.Warn("(fnModifyObject) params[1] >>>>>>>>>>", params[1])
		return params[0], nil
	}

	object := params[0].(map[string]interface{})
	key := params[1].(string)
	value := params[2]
	if nil != value {
		object[key] = value
	} else {
		delete(object, key)
	}

	log.Debug("(fnModifyObject) modified object = ", object)

	return object, nil
}
