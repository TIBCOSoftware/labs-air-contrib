package air

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnIsNumber{})
}

type fnIsNumber struct {
}

func (fnIsNumber) Name() string {
	return "isnumber"
}

func (fnIsNumber) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnIsNumber) Eval(params ...interface{}) (interface{}, error) {
	fmt.Println(params[0])

	switch dataType := params[0].(string); dataType {
	case "String":
		return false, nil
	case "Bool":
		return false, nil
	case "Binary":
		return false, nil

	case "Int8":
		return true, nil
	case "Int16":
		return true, nil
	case "Int32":
		return true, nil
	case "Int64":
		return true, nil
	case "Uint8":
		return true, nil
	case "Uint16":
		return true, nil
	case "Uint32":
		return true, nil
	case "Uint64":
		return true, nil
	case "Float32":
		return true, nil
	case "Float64":
		return true, nil
	}

	return nil, fmt.Errorf("Unknow data type %s", params[0])
}
