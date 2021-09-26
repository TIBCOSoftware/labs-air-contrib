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

	"github.com/SteveNY-Tibco/labs-air-contrib/activity/rules/internal/rules"
	"github.com/SteveNY-Tibco/labs-air-contrib/common/notification/notificationbroker"
	kwr "github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/keywordreplace"
	"github.com/SteveNY-Tibco/labs-lightcrane-contrib/common/util"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("tibco-f1-Rules")

var initialized bool = false

const (
	sLeftToken    = "leftToken"
	sRightToken   = "rightToken"
	sVariablesDef = "variablesDef"
	sTargets      = "targets"
	iData         = "Data"
	iGateway      = "gateway"
	iDevice       = "device"
	iEventID      = "id"
	iName         = "name"
	iValue        = "value"
	oSuccess      = "Success"
)

type Rules struct {
	metadata *activity.Metadata
	mux      sync.Mutex
	engines  map[string]*rules.RuleEngine
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aRules := &Rules{
		metadata: metadata,
		engines:  make(map[string]*rules.RuleEngine),
	}

	return aRules
}

func (a *Rules) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *Rules) Eval(context activity.Context) (done bool, err error) {
	log.Info("[Rules:Eval] entering ........ data = ", context.GetInput(iData))

	data, ok := context.GetInput(iData).(map[string]interface{})
	if !ok {
		return false, errors.New("Invalid data ... ")
	}
	gateway, ok := data[iGateway].(string)
	if !ok {
		return false, errors.New("Invalid device ... ")
	}
	device, ok := data[iDevice].(string)
	if !ok {
		return false, errors.New("Invalid device ... ")
	}
	eventID, ok := data[iEventID].(string)
	if !ok {
		return false, errors.New("Invalid eventID ... ")
	}
	name, ok := data[iName].(string)
	if !ok {
		return false, errors.New("Invalid name ... ")
	}
	value := data[iValue].(string)
	if !ok {
		return false, errors.New("Invalid value ... ")
	}

	log.Info(fmt.Sprintf("Received event from eventID: %s gateway: %s device: %s instrument: %s value: %s", eventID, gateway, device, name, value))

	engine, err := a.getRuleEngine(context)
	if nil != err {
		return true, err
	}

	engine.GetOrCreateResourceTuple(gateway, device, name, value)
	err = engine.Assert(eventID, gateway, device, name, value)
	if nil != err {
		return true, err
	}

	context.SetOutput(oSuccess, true)

	log.Info("[Rules:Eval] exit ........ ")

	return true, nil
}

func (a *Rules) SendNotification(notifier string, notification map[string]interface{}) error {
	log.Info("(Rules.SendNotification) notifier : ", notifier, ", notification : ", notification)
	notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
	if nil != notificationBroker {
		notificationBroker.SendEvent(notification)
	}
	return nil
}

func (a *Rules) getRuleEngine(ctx activity.Context) (*rules.RuleEngine, error) {
	var err error
	myId := util.ActivityId(ctx)
	engine := a.engines[myId]
	if nil == engine {
		a.mux.Lock()
		defer a.mux.Unlock()
		engine = a.engines[myId]
		if nil == engine {
			engineID, _ := ctx.GetSetting("id")
			tupleDescriptor, _ := ctx.GetSetting("tupleDescriptor")
			log.Info(fmt.Sprintf("Tuple Descriptor: %s", tupleDescriptor))
			defalutRuleDescriptor, _ := ctx.GetSetting("defalutRuleDescriptor")
			log.Info(fmt.Sprintf("Defalut Rule Descriptor: %s", defalutRuleDescriptor))

			engine = &rules.RuleEngine{}
			engine.SetEngineID(engineID.(string))
			engine.SetSender(a)
			err = engine.CreateRuleSessionThenStart(tupleDescriptor.(string))
			if nil != err {
				return nil, err
			}
			log.Info("Engine created ........ ")
			if nil != defalutRuleDescriptor && "{}" != defalutRuleDescriptor {
				log.Info("Build rule definitione ........ ")
				ruleDef := rules.RuleDefStruct{}
				if err := json.Unmarshal([]byte(defalutRuleDescriptor.(string)), &ruleDef); err != nil {
					log.Error(err)
					return nil, errors.New("Processing config request ERROR\n")
				}
				log.Infof("Processing addRule - Raw Object in main: %+v\n", ruleDef)
				engine.AddRule(ruleDef)
				log.Info("Processing addRule done .... ")
			}

			a.engines[myId] = engine
		}
		log.Info("engine = ", engine)
	}
	return engine, err
}

func (a *Rules) getVariableMapper(ctx activity.Context) *kwr.KeywordMapper {
	lefttoken, exist := ctx.GetSetting(sLeftToken)
	if !exist {
		return nil
	}
	righttoken, exist := ctx.GetSetting(sRightToken)
	if !exist {
		return nil
	}
	return kwr.NewKeywordMapper("", lefttoken.(string), righttoken.(string))
}
