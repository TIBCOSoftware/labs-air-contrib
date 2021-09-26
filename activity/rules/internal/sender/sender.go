package sender

type Sender interface {
	SendNotification(notifier string, notification map[string]interface{}) error
}
