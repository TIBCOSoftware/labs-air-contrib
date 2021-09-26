package f1

import (
	kwr "github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/keywordreplace"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnKewordReplace{})
}

type fnKewordReplace struct {
}

func (fnKewordReplace) Name() string {
	return "keywordreplace"
}

func (fnKewordReplace) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString, data.TypeString, data.TypeObject}, false
}

func (fnKewordReplace) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] || nil == params[1] || nil == params[2] || nil == params[3] {
		return params[0], nil
	}

	mapper := kwr.NewKeywordMapper(params[0].(string), params[1].(string), params[2].(string))
	return mapper.Replace("", params[3].(map[string]interface{})), nil
}
