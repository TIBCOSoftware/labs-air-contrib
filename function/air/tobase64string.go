package air

import (
	b64 "encoding/base64"
	"fmt"

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
	return []data.Type{data.TypeString}, false
}

func (fnToBase64String) Eval(params ...interface{}) (interface{}, error) {

	log.Debug("(fnToBase64String.Eval) params[0] : ", params[0])

	str, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("Illegal parameter!")
	}
	if "" == str {
		return "", nil
	}
	b64Str := b64.StdEncoding.EncodeToString([]byte(str))
	log.Debug("(fnToBase64String.Eval) encoded string : ", b64Str)

	return b64Str, nil
}
