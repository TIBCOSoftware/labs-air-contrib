package air

/*
	This function doesn't depend on EdgeX schema
*/

import (
	"fmt"
	"strings"
	"time"

	//"github.com/TIBCOSoftware/labs-flogo-lib/notification/notificationbroker"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnNotify{})
}

type fnNotify struct {
}

func (fnNotify) Name() string {
	return "notify"
}

func (fnNotify) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeObject, data.TypeArray, data.TypeString, data.TypeArray, data.TypeString}, false
}

func (fnNotify) Eval(params ...interface{}) (interface{}, error) {
	gateway := params[0]
	reading := params[1].(map[string]interface{})
	enriched := params[2].([]interface{})
	matchings := params[4].([]interface{})
	notifier := params[5].(string)

	target, error := getTarget(reading, enriched, params[3])
	if nil != error {
		return nil, error
	}

	log.Info("(fnNotify.Eval) gateway : ", gateway)
	log.Info("(fnNotify.Eval) reading : ", reading)
	log.Info("(fnNotify.Eval) enriched : ", enriched)

	matched := false
	matchType := ""
	level := ""
	for _, value := range matchings {
		match := value.(map[string]interface{})
		matchType = match["type"].(string)
		level = match["value"].(string)
		if "contains" == matchType && strings.Contains(target, level) {
			matched = true
			break
		}
	}

	reading["origin"] = time.Now().UnixNano() / 1e6

	if matched {
		notification := map[string]interface{}{
			"gateway": gateway,
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
					"value":    "TextMatching: " + matchType,
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "description",
					"value":    target[9 : len(target)-3],
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "level",
					"value":    level,
				},
			},
		}

		log.Info("(fnNotify.Eval) notifier : ", notifier, ", notification : ", notification)
		//notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
		//if nil != notificationBroker {
		//	notificationBroker.SendEvent(notification)
		//}
		return notification, nil
	}
	return nil, nil
}

func getTarget(reading map[string]interface{}, origEnriched []interface{}, format interface{}) (string, error) {

	enriched := make(map[string]interface{})
	for _, element := range origEnriched {
		enrichedElement := element.(map[string]interface{})
		enriched[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
	}

	data := NewKeywordMapper("@", "@").Replace(
		format.(string),
		NewKeywordReplaceHandler(reading, enriched),
	)

	log.Debug("(notifier.getTarget) out data string : ", data)

	return data, nil
}
