/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package notificationbroker

import (
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
	once.Do(func() {
		instance = &NotificationBrokerFactory{exeEventBrokers: make(map[string]*NotificationBroker)}
	})
	return instance
}

func (this *NotificationBrokerFactory) GetNotificationBroker(serverId string) *NotificationBroker {
	log.Info("(NotificationBrokerFactory.GetNotificationBroker) Factory : ", instance)
	log.Info("(NotificationBrokerFactory.GetNotificationBroker) EventBrokers : ", this.exeEventBrokers)
	return this.exeEventBrokers[serverId]
}

func (this *NotificationBrokerFactory) CreateNotificationBroker(
	brokerID string,
	listener NotificationListener) (*NotificationBroker, error) {

	this.mux.Lock()
	defer this.mux.Unlock()
	broker := this.exeEventBrokers[brokerID]

	broker = &NotificationBroker{
		ID:       brokerID,
		listener: listener,
	}
	this.exeEventBrokers[brokerID] = broker
	log.Info("(NotificationBrokerFactory.CreateNotificationBroker) Factory : ", instance)
	log.Info("(NotificationBrokerFactory.CreateNotificationBroker) EventBrokers : ", this.exeEventBrokers)

	return broker, nil
}

type NotificationBroker struct {
	ID       string
	listener NotificationListener
}

func (this *NotificationBroker) Start() {
	log.Info("(NotificationBroker.Start) Start broker, NotificationBroker : ", this)
}

func (this *NotificationBroker) Stop() {
}

func (this *NotificationBroker) SendEvent(event map[string]interface{}) {
	this.listener.ProcessEvent(this.ID, event)
}
