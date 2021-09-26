package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirData2Collection{})
}

type fnAirData2Collection struct {
}

func (fnAirData2Collection) Name() string {
	return "airdata2collection"
}

func (fnAirData2Collection) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeArray}, false
}

func (fnAirData2Collection) Eval(params ...interface{}) (interface{}, error) {

	log.Info("(fnAirData2Collection:Eval) entering ........ ")
	defer log.Info("(fnAirData2Collection:Eval) exit ........ ")

	airData := params[0].(map[string]interface{})
	var dataCollection []interface{}
	if nil != params[1] {
		dataCollection = params[1].([]interface{})
	} else {
		dataCollection = make([]interface{}, 0)
	}

	log.Debug("\n\n\n (fnAirData2Collection:Eval) in airData : ", airData)
	log.Debug("(fnAirData2Collection:Eval) in dataCollection : ", dataCollection)

	dataCollection = append(dataCollection, map[string]interface{}{
		"producer": "f1",
		"name":     "gateway",
		"value":    airData["gateway"],
	})

	for name, value := range airData["reading"].(map[string]interface{}) {
		log.Debug("(fnAirData2Collection:Eval) new name : ", name)
		dataCollection = append(dataCollection, map[string]interface{}{
			"producer": "f1",
			"name":     name,
			"value":    value,
		})
	}

	log.Debug("(fnAirData2Collection:Eval) out dataCollection : ", dataCollection)

	return dataCollection, nil
}
