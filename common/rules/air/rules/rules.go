package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/rules/common/model"
)

var log = logger.GetLogger("tibco-f1-dynamic")

// conditionCtxStruct - structure use to pass context to conditions
type conditionCtxStruct struct {
	Device                         string
	Resource                       string
	CompareNewMetricToValue        bool
	CompareNewMetricToValueOp      string
	CompareNewMetricValue          string
	CompareNewMetricToLastMetric   bool
	CompareNewMetricToLastMetricOp string
	CompareLastMetricToValue       bool
	CompareLastMetricToValueOp     string
	CompareLastMetricValue         string
}

func UpdateCond(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	log.Debug(fmt.Sprintf("Condition Evaluated: [%s]-[%s]\n", ruleName, condName))

	condResult := false

	readingTuple := tuples["ReadingEvent"]
	resourceTuple := tuples["ResourceConcept"]

	if readingTuple == nil || resourceTuple == nil || ctx == nil {
		log.Error("Should not get a nil Reading tuple or no context in updateCond! This is an error")
		return false
	}

	strCtx := ctx.(string)
	condCtx := conditionCtxStruct{}

	if err := json.Unmarshal([]byte(strCtx), &condCtx); err != nil {
		fmt.Printf("Processing config request ERROR\n")
	}

	readingTupleDevice, _ := readingTuple.GetString("device")
	readingTupleResource, _ := readingTuple.GetString("resource")
	resourceTupleDevice, _ := resourceTuple.GetString("device")
	resourceTupleResource, _ := resourceTuple.GetString("resource")

	if readingTupleResource == condCtx.Resource && readingTupleDevice == condCtx.Device &&
		resourceTupleResource == condCtx.Resource && resourceTupleDevice == condCtx.Device {
		condResult = true
	}

	return condResult
}

func UpdateAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	log.Debug(fmt.Sprintf("Rule fired: [%s]\n", ruleName))

	readingTuple := tuples["ReadingEvent"]
	resourceTuple := tuples["ResourceConcept"].(model.MutableTuple)

	readingTupleDevice, _ := readingTuple.GetString("device")
	readingTupleResource, _ := readingTuple.GetString("resource")
	log.Debug(fmt.Sprintf("Updating Device: [%s] Resource: [%s]\n", readingTupleDevice, readingTupleResource))

	// Update Value
	rtValue, _ := readingTuple.GetString("value")
	resourceTuple.SetString(nil, "value", rtValue)
}

// Condition to compare previous stored in a concept with the new reading value
func CompareValuesCond(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	log.Debug(fmt.Sprintf("Condition Evaluated: [%s]-[%s]\n", ruleName, condName))

	if true == true {
		return true
	}

	condResult := false

	readingTuple := tuples["ReadingEvent"]
	resourceTuple := tuples["ResourceConcept"]

	if readingTuple == nil || resourceTuple == nil || ctx == nil {
		log.Error("Should not get a nil Reading tuple or no context in compareValuesCond! This is an error")
		return false
	}

	strCtx := ctx.(string)
	condCtx := conditionCtxStruct{}

	if err := json.Unmarshal([]byte(strCtx), &condCtx); err != nil {
		fmt.Printf("Processing config request ERROR\n")
	}

	readingTupleDevice, _ := readingTuple.GetString("device")
	readingTupleResource, _ := readingTuple.GetString("resource")
	resourceTupleDevice, _ := resourceTuple.GetString("device")
	resourceTupleResource, _ := resourceTuple.GetString("resource")

	if readingTupleResource == condCtx.Resource && readingTupleDevice == condCtx.Device &&
		resourceTupleResource == condCtx.Resource && resourceTupleDevice == condCtx.Device {

		// Determine type
		readingTupleValueStr, _ := readingTuple.GetString("value")
		resourceTupleValueStr, _ := resourceTuple.GetString("value")
		isNumeric := false

		readingTupleValue, err := strconv.ParseFloat(readingTupleValueStr, 64)
		if err == nil {
			isNumeric = true
		}

		if condCtx.CompareNewMetricToValue {

			if isNumeric {
				// compare numbers

				value, _ := strconv.ParseFloat(condCtx.CompareNewMetricValue, 64)

				switch condCtx.CompareNewMetricToValueOp {
				case ">=":
					condResult = readingTupleValue > value
				case ">":
					condResult = readingTupleValue >= value
				case "<=":
					condResult = readingTupleValue <= value
				case "<":
					condResult = readingTupleValue < value
				case "==":
					condResult = readingTupleValue == value
				case "!=":
					condResult = readingTupleValue != value
				default:
					condResult = false
				}
			} else {
				// Compare string

				switch condCtx.CompareNewMetricToValueOp {
				case ">=":
					condResult = readingTupleValueStr > condCtx.CompareNewMetricValue
				case ">":
					condResult = readingTupleValueStr >= condCtx.CompareNewMetricValue
				case "<=":
					condResult = readingTupleValueStr <= condCtx.CompareNewMetricValue
				case "<":
					condResult = readingTupleValueStr < condCtx.CompareNewMetricValue
				case "==":
					condResult = readingTupleValueStr == condCtx.CompareNewMetricValue
				case "!=":
					condResult = readingTupleValueStr != condCtx.CompareNewMetricValue
				default:
					condResult = false
				}

				// log.Info(fmt.Sprintf("Evaluated tuple value operator context value: [%s][%s][%s] = [%t]",
				// 	readingTupleValueStr, condCtx.CompareNewMetricToValueOp, condCtx.CompareNewMetricValue, condResult))
			}

		}

		if condResult && condCtx.CompareNewMetricToLastMetric {

			if isNumeric {
				// compare numbers

				resourceTupleValue, _ := strconv.ParseFloat(resourceTupleValueStr, 64)

				switch condCtx.CompareNewMetricToLastMetricOp {
				case ">=":
					condResult = readingTupleValue > resourceTupleValue
				case ">":
					condResult = readingTupleValue >= resourceTupleValue
				case "<=":
					condResult = readingTupleValue <= resourceTupleValue
				case "<":
					condResult = readingTupleValue < resourceTupleValue
				case "==":
					condResult = readingTupleValue == resourceTupleValue
				case "!=":
					condResult = readingTupleValue != resourceTupleValue
				default:
					condResult = false
				}

			} else {
				// Compare string

				switch condCtx.CompareNewMetricToLastMetricOp {
				case ">=":
					condResult = readingTupleValueStr > resourceTupleValueStr
				case ">":
					condResult = readingTupleValueStr >= resourceTupleValueStr
				case "<=":
					condResult = readingTupleValueStr <= resourceTupleValueStr
				case "<":
					condResult = readingTupleValueStr < resourceTupleValueStr
				case "==":
					condResult = readingTupleValueStr == resourceTupleValueStr
				case "!=":
					condResult = readingTupleValueStr != resourceTupleValueStr
				default:
					condResult = false
				}
			}

		}

		if condResult && condCtx.CompareLastMetricToValue {

			if isNumeric {
				// compare numbers

				resourceTupleValue, _ := strconv.ParseFloat(resourceTupleValueStr, 64)
				value, _ := strconv.ParseFloat(condCtx.CompareLastMetricValue, 64)

				switch condCtx.CompareLastMetricToValueOp {
				case ">=":
					condResult = resourceTupleValue > value
				case ">":
					condResult = resourceTupleValue >= value
				case "<=":
					condResult = resourceTupleValue <= value
				case "<":
					condResult = resourceTupleValue < value
				case "==":
					condResult = resourceTupleValue == value
				case "!=":
					condResult = resourceTupleValue != value
				default:
					condResult = false
				}
			} else {
				// Compare string

				switch condCtx.CompareLastMetricToValueOp {
				case ">=":
					condResult = resourceTupleValueStr > condCtx.CompareLastMetricValue
				case ">":
					condResult = resourceTupleValueStr >= condCtx.CompareLastMetricValue
				case "<=":
					condResult = resourceTupleValueStr <= condCtx.CompareLastMetricValue
				case "<":
					condResult = resourceTupleValueStr < condCtx.CompareLastMetricValue
				case "==":
					condResult = resourceTupleValueStr == condCtx.CompareLastMetricValue
				case "!=":
					condResult = resourceTupleValueStr != condCtx.CompareLastMetricValue
				default:
					condResult = false
				}

			}

		}

	}

	log.Debug(fmt.Sprintf("Condition Evaluated: [%s]-[%s]-[%t]\n", ruleName, condName, condResult))

	return condResult
}

// Send notification whenever a compare value condition is true
/*
func CompareValuesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	log.Debug(fmt.Sprintf("Rule fired: rule = [%s], tuple =[%v]\n", ruleName, tuples))

	readingTuple := tuples["ReadingEvent"]
	resourceTuple := tuples["ResourceConcept"]

	if readingTuple == nil || resourceTuple == nil {
		log.Info("Should not get a nil tuple in FilterCondition! This is an error")
		return
	}

	gateway, _ := readingTuple.GetString("gateway")
	value, _ := readingTuple.GetString("value")
	device, _ := readingTuple.GetString("device")
	resource, _ := readingTuple.GetString("resource")
	ts := time.Now().UnixNano() / 1e6
	uuid := strconv.FormatInt(ts, 10)

	log.Debug(fmt.Sprintf("Action for: device = [%s], resource = [%s], value = [%s] \n", device, resource, value))

	if this.sendCommand {
		commandContext := map[string]interface{}{
			"gateway": gateway,
			"reading": map[string]interface{}{
				"origin": ts,
				"id":     uuid,
				"device": device,
				"name":   resource,
				"value":  value,
			},
			"enriched": []interface{}{
				map[string]interface{}{
					"producer": "rule",
					"name":     "Notification",
					"value":    "Command",
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "actionValue",
					"value":    this.actionValue,
				},
			},
		}

		log.Info(fmt.Sprintf("(RuleEngine.compareValuesAction) Marshalled commandContext: %v \n", commandContext))

		// Send Command
		this.sender.SendNotification(this.engineID, commandContext)
	}

	if this.sendNotification {
		//notification := "{\"id\": \"@f1..id@\",\"origin\": \"@f1..origin@\",\"device\": \"@f1..device@\",\"name\": \"@f1..name@\",\"value\": \"@f1..value@\",\"source\": \"@rule..source@\",\"description\": \"@rule..description@\",\"level\": \"@rule..level@\"}"
		notificationContext := map[string]interface{}{
			"gateway": gateway,
			"reading": map[string]interface{}{
				"origin": ts,
				"id":     uuid,
				"device": device,
				"name":   resource,
				"value":  value,
			},
			"enriched": []interface{}{
				map[string]interface{}{
					"producer": "rule",
					"name":     "Notification",
					"value":    "Message",
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "source",
					"value":    "Flogo Rule: " + ruleName,
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "description",
					"value":    ruleCtx.(string),
				},
				map[string]interface{}{
					"producer": "rule",
					"name":     "level",
					"value":    "INFO",
				},
			},
		}

		log.Info(fmt.Sprintf("(RuleEngine.compareValuesAction) Marshalled notificationContext: %v \n", notificationContext))

		// Send Notification
		this.sender.SendNotification(this.engineID, notificationContext)
	}
}
*/
