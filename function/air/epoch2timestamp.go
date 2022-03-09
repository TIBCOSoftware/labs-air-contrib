package air

import (
	"time"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnEpoch2Timestamp{})
}

type fnEpoch2Timestamp struct {
}

func (fnEpoch2Timestamp) Name() string {
	return "epoch2timestamp"
}

func (fnEpoch2Timestamp) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeInt, data.TypeString}, false
}

func (fnEpoch2Timestamp) Eval(params ...interface{}) (interface{}, error) {
	/* quick and dirty : only for realtime scenario */
	log.Info("(fnEpoch2Timestamp.Eval) params[0] : ", params[0], ", params[1] : ", params[1])
	epoch := int64(params[0].(int))
	if epoch > 1000000000000000000 {
		epoch = epoch / 1000000000
	} else if epoch > 1000000000000 {
		epoch = epoch / 1000
	}

	return time.Unix(epoch, 0).Format(params[1].(string)), nil
}
