/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package error

import (
	//	"errors"
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-air-contrib/common/notification/notificationbroker"
)

var log = logger.GetLogger("tibco-f1-Error")

var initialized bool = false

const (
	iActivity = "Activity"
	iMessage  = "Message"
	iData     = "Data"
	iGateway  = "Gateway"
	iReading  = "Reading"
	iEnriched = "Enriched"
	oSuccess  = "Success"
)

type Error struct {
	metadata *activity.Metadata
	mux      sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	aError := &Error{
		metadata: metadata,
	}

	return aError
}

func (a *Error) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *Error) Eval(context activity.Context) (done bool, err error) {
	log.Info("(Eval) entering ........ ")

	activity, ok := context.GetInput(iActivity).(string)
	if !ok {
		//return false, errors.New("Invalid activity ... ")
		log.Warn("Invalid activity ... ")
	}
	message, ok := context.GetInput(iMessage).(string)
	if !ok {
		//return false, errors.New("Invalid message ... ")
		log.Warn("Invalid message ... ")
	}
	data, ok := context.GetInput(iData).(string)
	if !ok {
		//return false, errors.New("Invalid data ... ")
		log.Warn("Invalid data ... ")
	}
	gateway, ok := context.GetInput(iGateway).(string)
	if !ok {
		//return false, errors.New("Invalid gateway ... ")
		log.Warn("Invalid gateway ... ")
	}

	reading, ok := context.GetInput(iReading).(map[string]interface{})
	if !ok {
		//return false, errors.New("Invalid reading ... ")
		log.Warn("Invalid reading ... ")
	}
	enriched, ok := context.GetInput(iEnriched).([]interface{})
	if !ok {
		//return false, errors.New("Invalid enriched ... ")
		log.Warn("Invalid enriched ... ")
	}

	log.Info(fmt.Sprintf("Received event from activity: %s gateway: %s message: %s data: %s", activity, gateway, message, data))
	oEnriched := []interface{}{
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "Notification",
			"value":    "Error",
		},
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "source",
			"value":    "Failed component: " + activity,
		},
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "description",
			"value":    message,
		},
		map[string]interface{}{
			"producer": "ErrorHandler",
			"name":     "data",
			"value":    data,
		},
	}

	for index := range enriched {
		oEnriched = append(oEnriched, enriched[index])
	}

	a.SendNotification("ErrorHandler", map[string]interface{}{
		"gateway":  gateway,
		"reading":  reading,
		"enriched": oEnriched,
	})

	context.SetOutput(oSuccess, true)

	log.Info("(Eval) exit ........ ")

	return true, nil
}

func (a *Error) SendNotification(notifier string, notification map[string]interface{}) error {
	log.Info("(Error.SendNotification) notifier : ", notifier, ", notification : ", notification)
	notificationBroker := notificationbroker.GetFactory().GetNotificationBroker(notifier)
	if nil != notificationBroker {
		notificationBroker.SendEvent(notification)
	}
	return nil
}
