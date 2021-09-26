package f1

import (
	"fmt"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnRTSFString2Properties{})
}

type fnRTSFString2Properties struct {
}

func (fnRTSFString2Properties) Name() string {
	return "rtsfstr2properties"
}

func (fnRTSFString2Properties) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

/*
customer_id:joe5,
employee_id:mary1,
event_time:1.573601393e+13,
lane_id:1,
basket_id:abc-012345-def
*/

func (fnRTSFString2Properties) Eval(params ...interface{}) (interface{}, error) {
	fmt.Println("params[0] = ", params[0])
	rtsfStr, ok1 := params[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("Illegal parameter : RTSF string")
	}
	fmt.Println("rtsfStr = ", rtsfStr)

	propertyArray := strings.Split(rtsfStr, ",")
	fmt.Println("propertyArray = ", propertyArray)
	properties := make(map[string]interface{})
	for _, propertyStr := range propertyArray {
		fmt.Println("propertyStr = ", propertyStr)
		if "" == propertyStr {
			continue
		}
		propertyElements := strings.Split(propertyStr, ":")
		properties[propertyElements[0]] = propertyElements[1]
	}

	return properties, nil
}
