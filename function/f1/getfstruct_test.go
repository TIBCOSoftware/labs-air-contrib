package f1

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/data/expression/function"
	"github.com/stretchr/testify/assert"
)

func TestFnLen_Eval(t *testing.T) {
	f := &fnGetFStruct{}
	v, err := function.Eval(f, "/Users/steven/Desktop/test/projects", 0)
	assert.Nil(t, err)
	fmt.Println("#### ", v)
}
