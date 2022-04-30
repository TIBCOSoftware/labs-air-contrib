package air

/*
	This function doesn't depend on EdgeX schema
*/

import (
	"fmt"
	"strings"
	"time"

	"github.com/TIBCOSoftware/labs-air-contrib/common/notification/notificationbroker"
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

	template := params[3].(string)
	iTarget, error := getTarget(reading, enriched, template)
	if nil != error {
		return nil, error
	}

	target := iTarget.(string)

	log.Debug("(fnNotify.Eval) gateway : ", gateway)
	log.Debug("(fnNotify.Eval) reading : ", reading)
	log.Debug("(fnNotify.Eval) enriched : ", enriched)
	log.Debug("(fnNotify.Eval) target : ", target)

	matched := false
	matchType := ""
	level := ""
	for _, value := range matchings {
		match := value.(map[string]interface{})
		matchType = match["type"].(string)
		level = match["value"].(string)
		log.Debug("(fnNotify.Eval) matchType : ", matchType, ", level : ", level, ", target : ", target, ", \"contains\" == matchType : ", "contains" == matchType, ", strings.Contains(target, level) : ", strings.Contains(target, level))
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

		log.Debug("(fnNotify.Eval) notifier : ", notifier, ", notification : ", notification)
		notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
		if nil != notificationBroker {
			go notificationBroker.SendEvent(notification)
		}
		return notification, nil
	}
	return nil, nil
}

func getTarget(reading map[string]interface{}, origEnriched []interface{}, template string) (interface{}, error) {

	dataMap := make(map[string]interface{})
	for key, value := range reading { // reading
		dataMap[fmt.Sprintf("f1..%s", key)] = value
	}
	for _, element := range origEnriched {
		enrichedElement := element.(map[string]interface{})
		dataMap[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
	}

	log.Debug("(getTarget) input dataMap : ", dataMap)
	log.Debug("(getTarget) input template : ", template)

	data := dataMap[template[1:len(template)-1]]
	if nil == data {
		data = NewKeywordMapper("@", "@").Replace(
			template,
			NewKeywordReplaceHandler(dataMap),
		)
	}
	log.Debug("(notifier.getTarget) out data : ", data)

	return data, nil
}
