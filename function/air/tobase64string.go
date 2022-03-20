package air

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnToBase64String{})
}

type fnToBase64String struct {
}

func (fnToBase64String) Name() string {
	return "tobase64string"
}

func (fnToBase64String) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeAny}, false
}

func (fnToBase64String) Eval(params ...interface{}) (interface{}, error) {

	log.Debug("(fnToBase64String.Eval) params[0] : ", params[0])

	if nil == params[0] {
		return nil, fmt.Errorf("(fnToBase64String.Eval) Illegal data : nil")
	}

	var data []byte
	if sData, ok := params[0].(string); ok {
		data = []byte(sData)
	} else if mData, ok := params[0].(map[string]interface{}); ok {
		jsonBuf, _ := json.Marshal(mData)
		data = jsonBuf
	} else if aData, ok := params[0].([]interface{}); ok {
		jsonBuf, _ := json.Marshal(aData)
		data = jsonBuf
	} else if bData, ok := params[0].([]byte); ok {
		data = bData
	} else {
		return nil, fmt.Errorf("(fnToBase64String.Eval) Illegal data type : %s", reflect.ValueOf(data).Kind().String())
	}

	b64Str := b64.StdEncoding.EncodeToString(data)
	log.Debug("(fnToBase64String.Eval) encoded string : ", b64Str)

	return b64Str, nil
}
