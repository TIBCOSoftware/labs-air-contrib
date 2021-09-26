/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package notificationlistener

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/SteveNY-Tibco/labs-air-contrib/common/notification/notificationbroker"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

const (
	cNotifierID = "notifierID"
)

//-============================================-//
//   Entry point register Trigger & factory
//-============================================-//

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&NotificationListener{}, &Factory{})
}

//-===============================-//
//     Define Trigger Factory
//-===============================-//

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(config.Settings, settings, true)
	if err != nil {
		return nil, err
	}

	return &NotificationListener{settings: settings}, nil
}

//-=========================-//
//      Define Trigger
//-=========================-//

var logger log.Logger

type NotificationListener struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	server   *notificationbroker.NotificationBroker
	mux      sync.Mutex

	settings  *Settings
	handlers  []trigger.Handler
	brokers   map[string]*notificationbroker.NotificationBroker
	listeners map[string](map[string]trigger.Handler)
}

// implements trigger.Initializable.Initialize
func (this *NotificationListener) Initialize(ctx trigger.InitContext) error {

	this.handlers = ctx.GetHandlers()
	logger = ctx.Logger()

	return nil
}

// implements ext.Trigger.Start
func (this *NotificationListener) Start() error {

	logger.Info("(NotificationListener.Start) Entering ... ")

	this.brokers = make(map[string]*notificationbroker.NotificationBroker)
	logger.Info("(NotificationListener.Start) this.brokers = ", this.brokers)
	this.listeners = make(map[string](map[string]trigger.Handler))
	for index, handler := range this.handlers {
		logger.Info("(NotificationListener.Start) handler = ", handler)
		broker, exist := handler.Settings()[cNotifierID]
		if !exist {
			return fmt.Errorf("Illegal broker : key = %s\n", cNotifierID)
		}

		brokerID := broker.(string)
		if nil == this.brokers[brokerID] {
			notifier, err := notificationbroker.GetFactory().CreateNotificationBroker(brokerID, this)
			if nil != err {
				return err
			}
			logger.Info("(NotificationListener.Start) Create Notifier = ", *notifier)
			this.brokers[brokerID] = notifier
			go notifier.Start()
		}
		handlers := this.listeners[brokerID]
		if nil == handlers {
			handlers = make(map[string]trigger.Handler)
			this.listeners[brokerID] = handlers
		}
		handlers[strconv.Itoa(index)] = handler
	}

	return nil
}

// implements ext.Trigger.Stop
func (this *NotificationListener) Stop() error {
	for _, notifier := range this.brokers {
		notifier.Stop()
	}
	return nil
}

func (this *NotificationListener) ProcessEvent(brokerID string, event map[string]interface{}) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	logger.Info("(NotificationListener.ProcessEvent) Got notification : id = ", brokerID, ", data = ", event)
	outputData := &Output{}
	outputData.Gateway = event["gateway"].(string)
	outputData.Reading = event["reading"].(map[string]interface{})
	outputData.Enriched = event["enriched"].([]interface{})

	for _, handler := range this.listeners[brokerID] {
		logger.Info("(NotificationListener.ProcessEvent) Send notification to flow : handler = ", handler.Name(), ", data = ", outputData)
		_, err := handler.Handle(context.Background(), outputData)
		if nil != err {
			logger.Info("(NotificationListener.ProcessEvent) Error -> ", err)
			return err
		}
	}

	return nil
}
