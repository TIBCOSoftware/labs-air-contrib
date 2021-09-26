package rules

import (
	"context"
	"encoding/json"
	"fmt"

	"strconv"
	"time"

	"github.com/SteveNY-Tibco/labs-air-contrib/activity/rules/internal/sender"
	"github.com/SteveNY-Tibco/labs-air-contrib/common/rules/air/rules"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/ruleapi"
)

var log = logger.GetLogger("tibco-f1-RuleEngine")

type ledBody struct {
	Led string
}

// RuleDefStruct - rule definition data
type RuleDefStruct struct {
	Name                               string `json:"name"`
	Description                        string `json:"description"`
	CondDevice                         string `json:"condDevice"`
	CondResource                       string `json:"condResource"`
	CondCompareNewMetricToValue        bool   `json:"condCompareNewMetricToValue"`
	CondCompareNewMetricToValueOp      string `json:"condCompareNewMetricToValueOp"`
	CondCompareNewMetricValue          string `json:"condCompareNewMetricValue"`
	CondCompareNewMetricToLastMetric   bool   `json:"condCompareNewMetricToLastMetric"`
	CondCompareNewMetricToLastMetricOp string `json:"condCompareNewMetricToLastMetricOp"`
	CondCompareLastMetricToValue       bool   `json:"condCompareLastMetricToValue"`
	CondCompareLastMetricToValueOp     string `json:"condCompareLastMetricToValueOp"`
	CondCompareLastMetricValue         string `json:"condCompareLastMetricValue"`
	ActionSendNotification             bool   `json:"actionSendNotification"`
	ActionNotification                 string `json:"actionNotification"`
	ActionSendCommand                  bool   `json:"actionSendCommand"`
	ActionDevice                       string `json:"actionDevice"`
	ActionResource                     string `json:"actionResource"`
	ActionValue                        string `json:"actionValue"`
}

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

type notificationCtxStruct struct {
	Created     int64  `json:"created"`
	UUID        string `json:"uuid"`
	Source      string `json:"source"`
	Gateway     string `json:"gateway"`
	Device      string `json:"device"`
	Resource    string `json:"resource"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

type caseCtxStruct struct {
	Description  string `json:"Description"`
	Level        string `json:"Level"`
	Notification string `json:"Notification"`
	Device       string `json:"device"`
	Origin       int64  `json:"origin"`
	DeviceId     string `json:"deviceId"`
}

type RuleEngine struct {
	//	gatewayId   string
	ruleSession      model.RuleSession
	engineID         string
	sendCommand      bool
	sendNotification bool
	sender           sender.Sender
	actionValue      string
}

func (this *RuleEngine) Assert(eventID string, gateway string, device string, name string, value string) error {
	// Assert Reading event
	tcl, _ := model.NewTupleWithKeyValues("ReadingEvent", eventID)
	tcl.SetString(nil, "gateway", gateway)
	tcl.SetString(nil, "device", device)
	tcl.SetString(nil, "resource", name)
	tcl.SetString(nil, "value", value)
	return this.ruleSession.Assert(nil, tcl)
}

// Set Sender
func (this *RuleEngine) SetSender(sender sender.Sender) {
	this.sender = sender
}

// Set Engine
func (this *RuleEngine) SetEngineID(engineID string) {
	this.engineID = engineID
}

// AddRule - add rule
func (this *RuleEngine) AddRule(ruleDef RuleDefStruct) {
	//fmt.Printf("Inside AddRule\n")
	//fmt.Printf("Raw Object: %+v\n", ruleDef)
	//fmt.Printf("Rule struct request device: %s\n", ruleDef.CondDevice)
	//fmt.Printf("Rule struct request resource: %s\n", ruleDef.CondResource)

	this.sendCommand = ruleDef.ActionSendCommand
	this.sendNotification = ruleDef.ActionSendNotification
	this.actionValue = ruleDef.ActionValue

	condContext := conditionCtxStruct{
		Device:                         ruleDef.CondDevice,
		Resource:                       ruleDef.CondResource,
		CompareNewMetricToValue:        ruleDef.CondCompareNewMetricToValue,
		CompareNewMetricToValueOp:      ruleDef.CondCompareNewMetricToValueOp,
		CompareNewMetricValue:          ruleDef.CondCompareNewMetricValue,
		CompareNewMetricToLastMetric:   ruleDef.CondCompareNewMetricToLastMetric,
		CompareNewMetricToLastMetricOp: ruleDef.CondCompareNewMetricToLastMetricOp,
		CompareLastMetricToValue:       ruleDef.CondCompareLastMetricToValue,
		CompareLastMetricToValueOp:     ruleDef.CondCompareLastMetricToValueOp,
		CompareLastMetricValue:         ruleDef.CondCompareLastMetricValue,
	}

	var condContextJSON []byte
	condContextJSON, err := json.Marshal(condContext)

	if err != nil {
		fmt.Printf("Rule request ERROR\n")
	}

	//fmt.Printf("Marshalled condContext: %s\n", string(condContextJSON))

	rule := ruleapi.NewRule(ruleDef.Name)
	rule.AddCondition("compareValuesCond", []string{"ReadingEvent", "ResourceConcept"}, this.compareValuesCond, string(condContextJSON))
	rule.SetAction(this.compareValuesAction)
	rule.SetContext(ruleDef.ActionNotification)
	rule.SetPriority(1)

	err = this.ruleSession.AddRule(rule)

	if err != nil {
		fmt.Printf("ERROR adding rule: %s\n", rule.GetName())
	} else {
		fmt.Printf("Added rule success for: %s\n", rule.GetName())
	}
	this.ruleSession.ReplayTuplesForRule(ruleDef.Name)

	updateRuleName := "Update" + ruleDef.Name
	rule1 := ruleapi.NewRule(updateRuleName)
	rule1.AddCondition("updateCond", []string{"ReadingEvent", "ResourceConcept"}, this.updateCond, string(condContextJSON))
	rule1.SetAction(this.updateAction)
	rule1.SetContext("Update")
	rule1.SetPriority(9)
	err = this.ruleSession.AddRule(rule1)

	if err != nil {
		fmt.Printf("ERROR adding update rule: %s\n", rule1.GetName())
	} else {
		fmt.Printf("Added rule success for: %s\n", rule1.GetName())
	}
	this.ruleSession.ReplayTuplesForRule(updateRuleName)

}

// DeleteRule - remove rule
func (this *RuleEngine) DeleteRule(ruleName string) {
	//fmt.Printf("Inside DeleteRule\n")

	this.ruleSession.DeleteRule(ruleName)

	updateRuleName := "Update" + ruleName
	this.ruleSession.DeleteRule(updateRuleName)

}

// RegisterConditionsAndActions - register rule conditions and actions
func (this *RuleEngine) RegisterConditionsAndActions() {

	log.Debug(fmt.Sprintf("Register Conditions and actions entering ... \n"))

	config.RegisterConditionEvaluator("updateCond", this.updateCond)
	config.RegisterActionFunction("updateAction", this.updateAction)

	config.RegisterConditionEvaluator("compareValuesCond", this.compareValuesCond)
	config.RegisterActionFunction("compareValuesAction", this.compareValuesAction)

	log.Debug(fmt.Sprintf("Register Conditions and actions done ... \n"))

}

// GetOrCreateResourceTuple - Gets or Creates an assertedtuple (resources are stateful concepts)
func (this *RuleEngine) GetOrCreateResourceTuple(gateway, device, resource, value string) model.MutableTuple {

	log.Debug(fmt.Sprintf("In GetOrCreateResourceTuple: [%s]-[%s]\n", device, resource))

	// Check if tuple already asserted
	key := device + "_" + resource
	tupleType := model.TupleType("ResourceConcept")
	tk, err := model.NewTupleKeyWithKeyValues(tupleType, key)

	if err != nil {
		log.Error(fmt.Sprintf("NewTupleKeyWithKeyValues failed for device-resource: [%s]-[%s]\n", device, resource))
	}

	log.Debug(fmt.Sprintf("In GetOrCreateResourceTuple Keys created: [%s]\n", tk.String()))

	conceptOld := this.ruleSession.GetAssertedTuple(tk)

	if conceptOld == nil {
		log.Debug(fmt.Sprintf("No concept found for: [%s]-[%s]\n", device, resource))
		concept, cerr := model.NewTupleWithKeyValues(tupleType, key)

		if cerr != nil {
			log.Error(fmt.Sprintf("Creating failed for device-resource: [%s]-[%s] - %s\n", device, resource, cerr))
		}

		concept.SetString(nil, "id", key)
		concept.SetString(nil, "gateway", gateway)
		concept.SetString(nil, "device", device)
		concept.SetString(nil, "resource", resource)
		concept.SetString(nil, "value", value)
		err := this.ruleSession.Assert(nil, concept)

		if err != nil {
			log.Error(fmt.Sprintf("Assert failed for device-resource: [%s]-[%s] - %s\n", device, resource, err))
		}

		return concept
	} else {
		log.Debug(fmt.Sprintf("Concept found for: [%s]-[%s]\n", device, resource))
		concept := conceptOld.(model.MutableTuple)

		return concept
	}

}

// Check if a concept needs to be updated with the latest reading
func (this *RuleEngine) updateCond(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return rules.UpdateCond(ruleName, condName, tuples, ctx)
}

// Update a concept with the latest reading
func (this *RuleEngine) updateAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	rules.UpdateAction(ctx, rs, ruleName, tuples, ruleCtx)
}

// Condition to compare previous stored in a concept with the new reading value
func (this *RuleEngine) compareValuesCond(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return rules.CompareValuesCond(ruleName, condName, tuples, ctx)
}

// Send notification whenever a compare value condition is true
func (this *RuleEngine) compareValuesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
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
