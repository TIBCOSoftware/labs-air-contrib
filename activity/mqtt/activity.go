package mqtt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/ssl"
)

var clients = make(map[string]mqtt.Client)
var clientMapMutex = &sync.Mutex{}

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
		client:   mqttClient,
		settings: settings,
		topic:    ParseTopic(settings.Topic),
	}
	return act, nil
}

type Activity struct {
	settings *Settings
	client   mqtt.Client
	topic    Topic
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

	topic := a.settings.Topic
	if params := input.TopicParams; len(params) > 0 {
		topic = a.topic.String(params)
	}
	ctx.Logger().Infof("MQTT client publishing, client id : %s, msg : %s", a.settings.Id, input.Message)
	if token := a.client.Publish(topic, byte(a.settings.Qos), a.settings.Retain, input.Message); token.Wait() && token.Error() != nil {
		ctx.Logger().Errorf("Error in publishing: %v", err)
		return true, token.Error()
	}

	ctx.Logger().Debugf("Published Message: %v", input.Message)

	return true, nil
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
