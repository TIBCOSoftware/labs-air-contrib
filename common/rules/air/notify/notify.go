package notify

import (
	"math"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("tibco-f1-dynamic")

func Eval(oldData map[string]interface{}, newData map[string]interface{}, thresholdArr interface{}) (string, string, string, bool) {
	log.Info("(notify.Eval) oldData : ", oldData)
	log.Info("(notify.Eval) newData : ", newData)
	log.Info("(notify.Eval) thresholdArr : ", thresholdArr)
	if 0 == len(oldData) {
		return "", "", "", false
	}

	oldReading := oldData["reading"].(map[string]interface{})
	newReading := newData["reading"].(map[string]interface{})
	log.Info("(notify.Eval) oldReading[\"value\"] : ", oldReading["value"])
	log.Info("(notify.Eval) newReading[\"value\"] : ", newReading["value"])

	newTriplet, err := triplet(newReading["value"].(string))
	oldTriplet, err := triplet(oldReading["value"].(string))
	if nil != err {
		log.Error("Error", err.Error())
	}

	source := ""
	notification := ""
	notificationLevel := ""
	thresholdf64 := make([]float64, 3)
	for _, th := range thresholdArr.([]interface{}) {
		thresholdValue := th.(map[string]interface{})
		source = newReading["name"].(string)
		log.Info("(notify.Eval) thresholdValue : ", thresholdValue)
		if thresholdValue["name"] == source {
			notification = thresholdValue["notification"].(string)
			notificationLevel = thresholdValue["notificationLevel"].(string)
			for index, value := range thresholdValue["value"].([]interface{}) {
				thresholdf64[index] = value.(float64)
			}
		}
	}

	if math.Abs(norm(newTriplet)-norm(oldTriplet)) > norm(thresholdf64) {
		return source, notification, notificationLevel, true
	}
	return "", "", "", false
}

func triplet(value string) ([]float64, error) {
	values := strings.Split(value, ",")
	e1, err := strconv.ParseFloat(values[0], 64)
	var e2 float64
	if len(values) >= 2 {
		e2, err = strconv.ParseFloat(values[1], 64)
	} else {
		e2 = float64(0)
	}
	var e3 float64
	if len(values) >= 3 {
		e3, err = strconv.ParseFloat(values[2], 64)
	} else {
		e3 = float64(0)
	}
	if nil != err {
		return nil, err
	}
	t := []float64{e1, e2, e3}
	return t, nil
}

func norm(vec []float64) float64 {
	return math.Sqrt(math.Pow(vec[0], 2) + math.Pow(vec[1], 2) + math.Pow(vec[2], 2))
}
