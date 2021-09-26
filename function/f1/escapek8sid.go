package f1

import (
	"regexp"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

var log = logger.GetLogger("tibco-f1_functions")

func init() {
	function.Register(&fnEscapeK8sID{})
}

type fnEscapeK8sID struct {
}

func (fnEscapeK8sID) Name() string {
	return "escapek8sid"
}

func (fnEscapeK8sID) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnEscapeK8sID) Eval(params ...interface{}) (interface{}, error) {
	if nil == params[0] {
		return params[0], nil
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Error(err)
	}
	return reg.ReplaceAllString(strings.ToLower(params[0].(string)), "-"), nil
}
