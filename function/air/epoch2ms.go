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
	/* quick and dirty : not for simulate future events */
	log.Info("(fnEpoch2ms.Eval) params[0] : ", params[0])
	epoch := int64(params[0].(int))
	if epoch > 9000000000000000 { // must be nano second - If it's in micro sec, it's GMT: Wednesday, March 14, 2255 4:00:00 PM
		epoch = epoch / 1000000
	} else if epoch > 10000000000000 { // must be micro second - If it's in millisec, it's GMT: Saturday, November 20, 2286 5:46:40 PM
		epoch = epoch / 1000
	}

	return epoch, nil
}
