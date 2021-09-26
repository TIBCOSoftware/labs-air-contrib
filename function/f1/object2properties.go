package f1

import (
	"strings"

	"github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/util"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnObject2Properties{})
}

type fnObject2Properties struct {
}

func (fnObject2Properties) Name() string {
	return "object2properties"
}

func (fnObject2Properties) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeString}, false
}

func (fnObject2Properties) Eval(params ...interface{}) (interface{}, error) {
	log.Info("(fnObject2Properties.Eval) params[0] : ", params[0], ", params[1]", params[1])
	object := params[0]
	properties := make([]interface{}, 0)
	if nil == object {
		return properties, nil
	}

	keys := params[1]
	objs := object.(map[string]interface{})
	if nil == keys || "" == keys.(string) {
		for key, value := range objs {
			sValue, err := util.ConvertToString(value, "")
			if nil != err {
				log.Warn("(fnObject2Properties.Eval) error parsing data : ", err.Error())
				continue
			}
			properties = append(properties, map[string]interface{}{
				"Name":  key,
				"Value": sValue,
			})
		}
	} else {
		keysArray := strings.Split(keys.(string), ",")
		for _, key := range keysArray {
			sValue, err := util.ConvertToString(objs[key], "")
			if nil != err {
				log.Warn("(fnObject2Properties.Eval) error parsing data : ", err.Error())
				continue
			}
			properties = append(properties, map[string]interface{}{
				"Name":  key,
				"Value": sValue,
			})
		}
	}

	return properties, nil
}
