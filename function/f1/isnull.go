package f1

import (
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnIsNull{})
}

type fnIsNull struct {
}

func (fnIsNull) Name() string {
	return "isnull"
}

func (fnIsNull) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject}, false
}

func (fnIsNull) Eval(params ...interface{}) (interface{}, error) {
	fmt.Println(params[0])
	return (nil == params[0]), nil
}
