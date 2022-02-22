package air

/*
	!!! This function DO depend on EdgeX schema !!!
*/

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirData2Array{})
}

type fnAirData2Array struct {
}

func (fnAirData2Array) Name() string {
	return "airdata2array"
}

func (fnAirData2Array) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeObject, data.TypeArray}, false
}

func (fnAirData2Array) Eval(params ...interface{}) (interface{}, error) {
	gateway := params[0]
	reading := params[1].(map[string]interface{})
	enriched := params[2].([]interface{})

	readingArray := make([]interface{}, 0)
	readingArray = append(readingArray, reading)
	for _, value := range enriched {
		log.Debug("(fnAirData2Array.Eval) new value : ", value)
		enriched := value.(map[string]interface{})
		newId := fmt.Sprintf("%s_%s", reading["id"], enriched["name"])
		if nil != enriched["producer"] && "" != enriched["producer"] {
			newId = fmt.Sprintf("%s_%s:%s", reading["id"], enriched["producer"], enriched["name"])
		}
		readingArray = append(readingArray, map[string]interface{}{
			"id":           newId,
			"origin":       reading["origin"],
			"deviceName":   reading["deviceName"],
			"resourceName": fmt.Sprintf("%s_%s", reading["resourceName"], enriched["name"]),
			"value":        enriched["value"],
			"valueType":    enriched["valueType"],
			"mediaType":    reading["mediaType"],
		})
	}

	dataArray := map[string]interface{}{
		"gateway":  gateway,
		"readings": readingArray,
	}
	log.Debug("(fnAirData2Array.Eval) out dataArray : ", dataArray)
	return dataArray, nil
}
