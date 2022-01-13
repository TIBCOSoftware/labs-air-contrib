package air

/*
	This function doesn't depend on EdgeX schema
*/

import (
	"encoding/json"
	"time"

	"github.com/TIBCOSoftware/labs-air-contrib/common/notification/notificationbroker"
	"github.com/TIBCOSoftware/labs-air-contrib/common/rules/air/notify"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnNotifyasy{})
}

type fnNotifyasy struct {
}

func (fnNotifyasy) Name() string {
	return "notifyasy"
}

func (fnNotifyasy) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeObject, data.TypeObject, data.TypeString, data.TypeString}, false
}

//func Eval(oldData interface{}, newData interface{}) (string, string, string, bool) {
//	return "", "", "", false
//}

func (fnNotifyasy) Eval(params ...interface{}) (interface{}, error) {
	newData := params[0].(map[string]interface{})
	oldData := params[1].(map[string]interface{})
	notifier := params[3].(string)

	log.Info("(fnNotifyasy.Eval) newData : ", newData)
	log.Info("(fnNotifyasy.Eval) oldData : ", oldData)
	log.Info("(fnNotifyasy.Eval) notifier : ", notifier)

	reading := newData["reading"].(map[string]interface{})
	reading["origin"] = time.Now().UnixNano() / 1e6

	var threshold interface{}
	err := json.Unmarshal([]byte(params[2].(string)), &threshold)
	if nil != err {
		return nil, err
	}
	source, description, level, notify := notify.Eval(oldData, newData, threshold)
	if notify {
		notification := map[string]interface{}{
			"gateway": newData["gateway"],
			"reading": reading,
			"enriched": []interface{}{
				map[string]interface{}{
					"producer": "rule",
					"name":     "Notification",
					"value":    "Message",
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "source",
					"value":    source,
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "description",
					"value":    description,
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "level",
					"value":    level,
				},
			},
		}

		log.Info("(fnNotifyasy.Eval) notifier started : ", notifier, ", notification : ", notification)
		notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
		if nil != notificationBroker {
			go notificationBroker.SendEvent(notification)
		}
		log.Info("(fnNotifyasy.Eval) notifier done : ", notifier, ", notification : ", notification)
		return notification, nil
	}
	return nil, nil
}
