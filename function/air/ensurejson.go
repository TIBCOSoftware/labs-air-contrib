package air

import (
	//	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnEnsureJson{})
}

type fnEnsureJson struct {
}

func (fnEnsureJson) Name() string {
	return "ensurejson"
}

func (fnEnsureJson) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnEnsureJson) Eval(params ...interface{}) (interface{}, error) {
	log.Debug("(fnEnsureJson.Eval) params[0] : ", params[0])

	in, ok := params[0].(string)
	if !ok {
		err := fmt.Errorf("Illegal parameter!")
		log.Error("(fnEnsureJson.Eval) err : ", err.Error())
		return in, err
	}

	if false == strings.HasPrefix(in, "\"") {
		return in, nil
	}

	jsonStr, err := strconv.Unquote(in)
	if nil != err {
		log.Error("(fnEnsureJson.Eval) Unquote fail, err : ", err.Error())
		return in, err
	}

	return jsonStr, nil
}
