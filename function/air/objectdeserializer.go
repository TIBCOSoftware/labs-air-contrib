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
	function.Register(&fnObjectDeserializer{})
}

type fnObjectDeserializer struct {
}

func (fnObjectDeserializer) Name() string {
	return "objectdeserializer"
}

func (fnObjectDeserializer) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeBytes, data.TypeString}, false
}

func (fnObjectDeserializer) Eval(params ...interface{}) (interface{}, error) {

	log.Debug("(fnObjectDeserializer.Eval) params[0] : ", params[0])
	log.Debug("(fnObjectDeserializer.Eval) params[1] : ", params[1])

	objbytes, ok := params[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("Illegal object bytes!")
	}

	deserializer, ok := params[1].(string)
	if "" == deserializer || !ok {
		/* no predefined deserializer */
		deserializer = "JSON"
		if objbytes[0] != byte('{') && objbytes[0] != byte('[') {
			// If not JSON then assume it is CBOR
			deserializer = "CBOR"
		}
	}

	var obj interface{}
	switch deserializer {
	case "JSON":
		err := json.Unmarshal(objbytes, &obj)
		if err != nil {
			return nil, err
		}
	case "CBOR":
		//		eventRequest := &requests.AddEventRequest{}
		//		err := cbor.Unmarshal(objbytes, eventRequest)
		//		if err != nil {
		//			return nil, err
		//		}
		//		obj = eventRequest
		err := cbor.Unmarshal(objbytes, &obj)
		if err != nil {
			return nil, err
		}
	case "Base64JSON":
		decodedByteContent, err := base64.StdEncoding.DecodeString(string(objbytes))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(decodedByteContent, &obj)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Illegal deserializer : %s", deserializer)
	}

	log.Debug("(fnObjectDeserializer.Eval) deserializer : ", deserializer)
	log.Debug("(fnObjectDeserializer.Eval) obj : ", obj)

	return obj, nil
}
