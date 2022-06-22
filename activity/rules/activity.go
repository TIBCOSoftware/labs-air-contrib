/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/labs-air-contrib/activity/rules/internal/rules"
	"github.com/TIBCOSoftware/labs-air-contrib/common/notification/notificationbroker"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

const (
	iGateway = "gateway"
	iDevice  = "deviceName"
	iEventID = "id"
	iName    = "resourceName"
	iValue   = "value"
	oSuccess = "Success"
)

type Settings struct {
	ID                    string `md:"id"`
	TupleDescriptor       string `md:"tupleDescriptor"`
	DefalutRuleDescriptor string `md:"defalutRuleDescriptor"`
}

type Input struct {
	Data           map[string]interface{} `md:"Data"`
	RuleDescriptor string                 `md:"RuleDescriptor"`
}

type Output struct {
	Success bool `md:"Success"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Data":           i.Data,
		"RuleDescriptor": i.RuleDescriptor,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	ok := true
	i.Data, ok = values["Data"].(map[string]interface{})
	if !ok {
		return errors.New("Illegal Data type, expect map[string]interface{}.")
	}
	i.RuleDescriptor, ok = values["RuleDescriptor"].(string)
	if !ok {
		return errors.New("Illegal RuleDescriptor type, expect string.")
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Success": o.Success,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	o.Success = values["Success"].(bool)
	return nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

type Activity struct {
	mux    sync.Mutex
	engine *rules.RuleEngine
}

func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	engineID := settings.ID
	tupleDescriptor := settings.TupleDescriptor
	defalutRuleDescriptor := settings.DefalutRuleDescriptor

	engine := &rules.RuleEngine{}
	engine.SetEngineID(engineID)
	err = engine.CreateRuleSessionThenStart(tupleDescriptor)
	if nil != err {
		return nil, err
	}
	if "" != defalutRuleDescriptor && "{}" != defalutRuleDescriptor {
		ruleDef := rules.RuleDefStruct{}
		if err := json.Unmarshal([]byte(defalutRuleDescriptor), &ruleDef); err != nil {
			return nil, errors.New("Processing config request ERROR\n")
		}
		engine.AddRule(ruleDef)
	}

	activity := &Activity{
		engine: engine,
	}
	engine.SetSender(activity)

	return activity, nil
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	log := ctx.Logger()
	input := &Input{}
	ctx.GetInputObject(input)

	data := input.Data

	log.Info("(Eval) entering ........ data = ", data)
	defer log.Info("(Eval) exit ........ ")

	if nil == data {
		return false, errors.New("Invalid data ... ")
	}
	gateway, ok := data[iGateway].(string)
	if !ok {
		return false, errors.New("Invalid gateway ... ")
	}
	device, ok := data[iDevice].(string)
	if !ok {
		return false, errors.New("Invalid deviceName ... ")
	}
	eventID, ok := data[iEventID].(string)
	if !ok {
		return false, errors.New("Invalid eventID ... ")
	}
	name, ok := data[iName].(string)
	if !ok {
		return false, errors.New("Invalid resourceName ... ")
	}
	value := data[iValue].(string)
	if !ok {
		return false, errors.New("Invalid value ... ")
	}

	log.Info(fmt.Sprintf("(Eval) Received event from eventID: %s gateway: %s deviceName: %s instrument: %s value: %s", eventID, gateway, device, name, value))

	a.engine.GetOrCreateResourceTuple(gateway, device, name, value)
	err = a.engine.Assert(eventID, gateway, device, name, value)
	if nil != err {
		return true, err
	}

	ctx.SetOutput(oSuccess, true)

	return true, nil
}

func (a *Activity) SendNotification(notifier string, notification map[string]interface{}) error {
	notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
	if nil != notificationBroker {
		notificationBroker.SendEvent(notification)
	}
	return nil
}
