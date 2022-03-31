package air

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnEpoch2ms{})
}

type fnEpoch2ms struct {
}

func (fnEpoch2ms) Name() string {
	return "epoch2ms"
}

func (fnEpoch2ms) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeInt}, false
}

func (fnEpoch2ms) Eval(params ...interface{}) (interface{}, error) {
	/* quick and dirty : only for realtime scenario */
	log.Info("(fnEpoch2ms.Eval) params[0] : ", params[0])
	epoch := int64(params[0].(int))
	if epoch > 1000000000000000000 {
		epoch = epoch / 1000000000
	} else if epoch > 1000000000000 {
		epoch = epoch / 1000
	}

	return epoch, nil
}
