package air

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFnThrowException(t *testing.T) {
	f := &fnThrowException{}
	v, err := function.Eval(f, true, "Test Exception!!")
	assert.NotNil(t, err)
	fmt.Println("#### ", v)
	fmt.Println("#### ", err.Error())
}
