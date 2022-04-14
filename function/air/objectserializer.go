package air

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	//	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/fxamacker/cbor/v2"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnObjectSerializer{})
}

type fnObjectSerializer struct {
}

func (fnObjectSerializer) Name() string {
	return "objectserializer"
}

func (fnObjectSerializer) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny, data.TypeString}, false
}

func (fnObjectSerializer) Eval(params ...interface{}) (interface{}, error) {

	log.Info("(fnObjectSerializer.Eval) params[0] : ", params[0])
	log.Debug("(fnObjectSerializer.Eval) params[1] : ", params[1])

	if nil == params[0] {
		return nil, fmt.Errorf("Illegal nil object!")
	}

	deserializer, _ := params[1].(string)

	var objBytes []byte
	var err error
	switch deserializer {
	case "JSON":
		objBytes, err = json.Marshal(params[0])
		if err != nil {
			return nil, err
		}
	case "CBOR":
		objBytes, err = cbor.Marshal(params[0])
		if err != nil {
			return nil, err
		}
	case "Base64JSON":
		objBytes, err = json.Marshal(params[0])
		if err != nil {
			return nil, err
		}
		objBytes = []byte(base64.StdEncoding.EncodeToString(objBytes))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Illegal deserializer : %s", deserializer)
	}

	log.Debug("(fnObjectSerializer.Eval) deserializer : ", deserializer)
	log.Debug("(fnObjectSerializer.Eval) objBytes : ", objBytes)

	return objBytes, nil
}
