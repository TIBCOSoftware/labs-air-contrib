package mqttcoupler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/Joker666/AsyncGoDemo/async"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/ssl"
)

const (
	oData = "Data"
)

var clients = make(map[string]mqtt.Client)
var clientMapMutex = &sync.Mutex{}
var gotMessage = make(chan bool, 1)

func getClient(logger log.Logger, connectionID string, opts *mqtt.ClientOptions) (client mqtt.Client, err error) {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()

	client = clients[connectionID]
	if client != nil {
		logger.Debug("[mqtt.activity.getClient] Mqtt Publish is reusing an existing connection...")
		return client, nil
	}

	client = mqtt.NewClient(opts)
	logger.Info("[mqtt.activity.getClient] Mqtt Publish is establishing a connection, client id : ", connectionID)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, errors.New(fmt.Sprintf("Connection to mqtt broker failed %v", token.Error()))
	}

	clients[connectionID] = client

	return client, nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

// TokenType is a type of token
type TokenType int

const (
	// Literal is a literal token type
	Literal TokenType = iota
	// Substitution is a parameter substitution
	Substitution
)

// Token is a MQTT topic token
type Token struct {
	TokenType TokenType
	Token     string
}

// Topic is a parsed topic
type Topic []Token

// ParseTopic parses the topic
func ParseTopic(topic string) Topic {
	var parsed Topic
	parts, index := strings.Split(topic, "/"), 0
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			token := strings.TrimPrefix(part, ":")
			if token == "" {
				token = strconv.Itoa(index)
				index++
			}
			parsed = append(parsed, Token{
				TokenType: Substitution,
				Token:     token,
			})
		} else {
			parsed = append(parsed, Token{
				TokenType: Literal,
				Token:     part,
			})
		}
	}
	return parsed
}

// String generates a string for the topic with params
func (t Topic) String(params map[string]string) string {
	output := strings.Builder{}
	for i, token := range t {
		if i > 0 {
			output.WriteString("/")
		}
		switch token.TokenType {
		case Literal:
			output.WriteString(token.Token)
		case Substitution:
			if value, ok := params[token.Token]; ok {
				output.WriteString(value)
			} else {
				output.WriteString(":")
				output.WriteString(token.Token)
			}
		}
	}
	return output.String()
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	options := initClientOption(ctx.Logger(), settings)

	if strings.HasPrefix(settings.Broker, "ssl") {

		cfg := &ssl.Config{}

		if len(settings.SSLConfig) != 0 {
			err := cfg.FromMap(settings.SSLConfig)
			if err != nil {
				return nil, err
			}

			if _, set := settings.SSLConfig["skipVerify"]; !set {
				cfg.SkipVerify = true
			}
			if _, set := settings.SSLConfig["useSystemCert"]; !set {
				cfg.UseSystemCert = true
			}
		} else {
			//using ssl but not configured, use defaults
			cfg.SkipVerify = true
			cfg.UseSystemCert = true
		}

		tlsConfig, err := ssl.NewClientTLSConfig(cfg)
		if err != nil {
			return nil, err
		}

		options.SetTLSConfig(tlsConfig)
	}

	mqttClient, err := getClient(ctx.Logger(), settings.Id, options)
	if nil != err {
		return nil, err
	}

	act := &Activity{
		client:        mqttClient,
		settings:      settings,
		topic:         ParseTopic(settings.Topic),
		topicMessages: NewTopicMessages(),
	}
	return act, nil
}

type Activity struct {
	settings      *Settings
	client        mqtt.Client
	topic         Topic
	topicMessages *TopicMessages
	c             *sync.Cond
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}

	err = ctx.GetInputObject(input)

	if err != nil {
		return true, err
	}

	responseTimeout := a.settings.ResponseTimeout
	ctx.Logger().Infof("MQTT response timeout : %d", responseTimeout)

	topic := a.settings.Topic
	if params := input.TopicParams; len(params) > 0 {
		topic = a.topic.String(params)
	}

	uuid, _ := uuid.NewUUID()
	respTopic := fmt.Sprintf("%s/%s", topic, uuid)
	ctx.Logger().Infof("MQTT client sunscribe topic, client id : %s, topic : %s", a.settings.Id, respTopic)
	if token := a.client.Subscribe(respTopic, 0, a.topicMessages.receiveMessage); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return true, token.Error()
	}

	future := async.Exec(func() interface{} {
		return a.WaitMessage(ctx.Logger(), respTopic)
	})

	ctx.Logger().Infof("MQTT client publishing, client id : %s, msg : %s", a.settings.Id, input.Message)
	if token := a.client.Publish(topic, byte(a.settings.Qos), a.settings.Retain, input.Message); token.Wait() && token.Error() != nil {
		ctx.Logger().Errorf("Error in publishing: %v", err)
		return true, token.Error()
	}
	a.topicMessages.startTimer(ctx.Logger(), responseTimeout)
	ctx.Logger().Debugf("Request Message: %v", input.Message)

	msg := future.Await()
	ctx.Logger().Infof("Response Message: %v", msg)

	err = ctx.SetOutput(oData, msg)
	if err != nil {
		return false, err
	}

	_ = a.client.Unsubscribe(respTopic)

	return true, nil
}

func (a *Activity) WaitMessage(logger log.Logger, topic string) mqtt.Message {
	logger.Debugf("Awaiting message on : %s", topic)
	<-gotMessage
	msg := a.topicMessages.removeMessage(topic)
	logger.Debugf("Got signal! here is the message for the topic : %v", msg)
	return msg
}

func NewTopicMessages() *TopicMessages {
	topicMessages := &TopicMessages{
		messages: make(map[string]interface{}),
	}
	return topicMessages
}

type TopicMessages struct {
	activity *Activity
	messages map[string]interface{}
}

func (t *TopicMessages) startTimer(logger log.Logger, timeout int) {
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		logger.Debugf("Time up after sleep : %s seconds", timeout)
		gotMessage <- true
	}()
}

func (t *TopicMessages) receiveMessage(client mqtt.Client, msg mqtt.Message) {
	t.messages[msg.Topic()] = msg
	fmt.Println("xxxxxxxxxxxxxxxxxxx receiveMessage xxxxxxxxxxxxxxxxxxxxx Topic = ", msg.Topic())
	gotMessage <- true
}

func (t *TopicMessages) removeMessage(topic string) mqtt.Message {
	msg := t.messages[topic]
	if nil != msg {
		delete(t.messages, topic)
		return msg.(mqtt.Message)
	}
	return nil
}

func initClientOption(logger log.Logger, settings *Settings) *mqtt.ClientOptions {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(settings.Broker)
	opts.SetClientID(settings.Id)
	opts.SetUsername(settings.Username)
	password := settings.Password
	if strings.HasPrefix(password, "SECRET:") {
		pwdBytes, _ := base64.StdEncoding.DecodeString(password[7:])
		password = string(pwdBytes)
	}
	opts.SetPassword(password)
	opts.SetCleanSession(settings.CleanSession)
	if settings.KeepAlive != 0 {
		opts.SetKeepAlive(time.Duration(settings.KeepAlive) * time.Second)
	} else {
		opts.SetKeepAlive(2 * time.Second)
	}

	if settings.Store != "" && settings.Store != ":memory:" {
		logger.Debugf("Using file store: %s", settings.Store)
		opts.SetStore(mqtt.NewFileStore(settings.Store))
	}

	return opts
}
