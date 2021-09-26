package f1

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnCombineObjs{})
}

type fnCombineObjs struct {
}

func (fnCombineObjs) Name() string {
	return "combineobjs"
}

func (fnCombineObjs) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject}, true
}

func (fnCombineObjs) Eval(params ...interface{}) (interface{}, error) {
	combined := make(map[string]interface{})
	for _, param := range params {
		if nil != param {
			for key, value := range param.(map[string]interface{}) {
				combined[key] = value
			}
		}
	}

	return combined, nil
}
