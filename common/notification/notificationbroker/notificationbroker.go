/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package notificationbroker

import (
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("tibco-f1_notification_broker")

var (
	instance *NotificationBrokerFactory
	once     sync.Once
)

type NotificationListener interface {
	ProcessEvent(notifier string, event map[string]interface{}) error
}

type NotificationBrokerFactory struct {
	exeEventBrokers map[string]*NotificationBroker
	mux             sync.Mutex
}

func GetFactory() *NotificationBrokerFactory {
	log.Debug("(NotificationBrokerFactory.GetFactory) entering ..... ")
	once.Do(func() {
		instance = &NotificationBrokerFactory{exeEventBrokers: make(map[string]*NotificationBroker)}
	})
	log.Debug("(NotificationBrokerFactory.GetFactory) exit : Factory = ", instance)
	return instance
}

func (this *NotificationBrokerFactory) GetNotificationBroker(serverId string) *NotificationBroker {
	log.Debug("(NotificationBrokerFactory.GetNotificationBroker) Factory : ", instance)
	log.Debug("(NotificationBrokerFactory.GetNotificationBroker) EventBrokers : ", this.exeEventBrokers)
	return this.exeEventBrokers[serverId]
}

func (this *NotificationBrokerFactory) CreateNotificationBroker(
	brokerID string,
	listener NotificationListener) (*NotificationBroker, error) {

	log.Debug("(NotificationBrokerFactory.CreateNotificationBroker) Factory : ", instance)
	this.mux.Lock()
	defer this.mux.Unlock()
	broker := this.exeEventBrokers[brokerID]
	if nil == broker {
		broker = &NotificationBroker{
			ID:       brokerID,
			listener: listener,
		}
	}
	this.exeEventBrokers[brokerID] = broker
	log.Debug("(NotificationBrokerFactory.CreateNotificationBroker) EventBrokers : ", this.exeEventBrokers)

	return broker, nil
}

func (this *NotificationBrokerFactory) CreateNotificationBrokers(
	brokerIDsStr string,
	listener NotificationListener) ([]*NotificationBroker, error) {

	log.Debug("(NotificationBrokerFactory.CreateNotificationBrokers) Factory : ", instance)

	var brokers []*NotificationBroker
	brokerIds := strings.Split(brokerIDsStr, ",")
	for _, brokerID := range brokerIds {
		broker, err := this.CreateNotificationBroker(brokerID, listener)
		if nil != err {
			return nil, err
		}
		brokers = append(brokers, broker)
	}

	log.Debug("(NotificationBrokerFactory.CreateNotificationBrokers) EventBrokers : ", this.exeEventBrokers)

	return brokers, nil
}

type NotificationBroker struct {
	ID       string
	listener NotificationListener
}

func (this *NotificationBroker) Start() {
	log.Debug("(NotificationBroker.Start) Start broker, NotificationBroker : ", this)
}

func (this *NotificationBroker) Stop() {
}

func (this *NotificationBroker) SendEvent(event map[string]interface{}) {
	log.Debug("(NotificationBroker.SendEvent) event : ", event)
	this.listener.ProcessEvent(this.ID, event)
}
