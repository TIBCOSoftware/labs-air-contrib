package f1

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirData2ArrayOld{})
}

type fnAirData2ArrayOld struct {
}

func (fnAirData2ArrayOld) Name() string {
	return "airdata2arrayold"
}

func (fnAirData2ArrayOld) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeInt, data.TypeString, data.TypeString, data.TypeObject, data.TypeArray}, false
}

func (fnAirData2ArrayOld) Eval(params ...interface{}) (interface{}, error) {
	id := params[0]
	source := params[1]
	device := params[2]
	gateway := params[3]
	reading := params[4].(map[string]interface{})
	enriched := params[5].([]interface{})

	readingArray := make([]interface{}, 0)
	readingArray = append(readingArray, reading)
	for _, value := range enriched {
		log.Debug("(fnAirData2ArrayOld.Eval) new value ========>", value)
		enriched := value.(map[string]interface{})
		readingArray = append(readingArray, map[string]interface{}{
			"id":        fmt.Sprintf("%s_%s", reading["id"], enriched["name"]),
			"origin":    reading["origin"],
			"device":    reading["device"],
			"name":      fmt.Sprintf("%s_%s", reading["name"], enriched["name"]),
			"value":     enriched["value"],
			"valueType": enriched["type"],
			"mediaType": reading["mediaType"],
		})
	}

	dataArray := map[string]interface{}{
		"id":       id,
		"source":   source,
		"device":   device,
		"gateway":  gateway,
		"readings": readingArray,
	}
	log.Debug("(fnAirData2ArrayOld.Eval) out dataArray ========>", dataArray)
	return dataArray, nil
}
