/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package error

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/labs-air-contrib/common/notification/notificationbroker"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
)

const (
	EXTRA_DATA = "[[AIR::EXTRA_DATA]]"
)

func init() {
	_ = activity.Register(&Activity{})
}

type Input struct {
	Activity string                 `md:"Activity"`
	Message  string                 `md:"Message"`
	Data     string                 `md:"Data"`
	Gateway  string                 `md:"Gateway"`
	Reading  map[string]interface{} `md:"Reading"`
	Enriched []interface{}          `md:"Enriched"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Activity": i.Activity,
		"Message":  i.Message,
		"Data":     i.Data,
		"Gateway":  i.Gateway,
		"Reading":  i.Reading,
		"Enriched": i.Enriched,
	}
}

type Output struct {
	Success bool `md:"Success"`
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.Activity, err = coerce.ToString(values["Activity"])
	if err != nil {
		return err
	}
	i.Message, err = coerce.ToString(values["Message"])
	if err != nil {
		return err
	}
	i.Data, err = coerce.ToString(values["Data"])
	if err != nil {
		return err
	}
	i.Gateway, err = coerce.ToString(values["Gateway"])
	if err != nil {
		return err
	}
	i.Reading, err = coerce.ToObject(values["Reading"])
	if err != nil {
		return err
	}
	i.Enriched, err = coerce.ToArray(values["Enriched"])
	if err != nil {
		return err
	}

	return nil
}

var activityMd = activity.ToMetadata(&Input{})

// Activity is an Activity that is used to send error
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - send error
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	log := ctx.Logger()
	log.Info("(Eval) entering ........ ")

	input := &Input{}
	ctx.GetInputObject(input)

	log.Info(fmt.Sprintf("Received event input: %v", input))
	log.Info(fmt.Sprintf("Received event from activity: %s gateway: %s message: %s data: %s", input.Activity, input.Gateway, input.Message, input.Data))

	var result interface{}
	extraDataPos := strings.Index(input.Message, EXTRA_DATA)
	if 0 <= extraDataPos {
		extraDataStr := input.Message[extraDataPos+len(EXTRA_DATA):]
		if err := json.Unmarshal([]byte(extraDataStr), &result); nil != err {
			log.Warn("No result data back from failed component !")
		}
	}

	oEnriched := []interface{}{
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "ErrorCode",
			"value":    "300",
		},
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "source",
			"value":    input.Activity,
		},
		map[string]interface{}{
			"producer": input.Activity,
			"name":     "result",
			"value":    result,
		},
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "description",
			"value":    input.Message,
		},
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "data",
			"value":    input.Data,
		},
	}

	for index := range input.Enriched {
		oEnriched = append(oEnriched, input.Enriched[index])
	}

	a.SendNotification("ErrorHandler", map[string]interface{}{
		"gateway":  input.Gateway,
		"reading":  input.Reading,
		"enriched": oEnriched,
	})

	ctx.SetOutput("Success", true)

	log.Info("(Eval) exit ........ ")

	return true, nil
}

func (a *Activity) SendNotification(notifier string, notification map[string]interface{}) error {
	notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
	if nil != notificationBroker {
		notificationBroker.SendEvent(notification)
	}
	return nil
}
