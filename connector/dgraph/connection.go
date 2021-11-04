package dgraph

import (
	b64 "encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/TIBCOSoftware/labs-air-contrib/common/graphbuilder/dbservice"
	"github.com/TIBCOSoftware/labs-air-contrib/common/graphbuilder/dbservice/dgraph/services"
	"github.com/TIBCOSoftware/labs-air-contrib/common/graphbuilder/model"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "Dgraph.connection")
var factory = &DgraphFactory{}

func NewSetting(settings map[string]interface{}) (*Settings, error) {
	s := &Settings{}

	var err = metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}

	if s.Name == "" {
		return nil, errors.New("Required Parameter Name is missing")
	}

	//	cDescription := s.Description

	if s.ApiVersion == "" {
		return nil, errors.New("Required Parameter Model is missing")
	}

	if s.URL == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}

	if s.User == "" {
		logCache.Debug("Parameter User is empty")
	}

	if s.Password == "" {
		logCache.Debug("Parameter Password is empty")
	}

	if true == s.TLSEnabled && s.TLS == "" {
		return nil, errors.New("Required Parameter TLS is missing")
	}

	if s.SchemaGen == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}

	if s.Schema == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}
	return s, nil
}

// Settings for dgraph
type Settings struct {
	Name        string `md:"name,required"`
	Description string `md:"description"`
	ApiVersion  string `md:"apiVersion,required"`
	URL         string `md:"url"`
	TLSEnabled  bool   `md:"tlsEnabled,required"`
	User        string `md:"user,required"`
	Password    string `md:"password,required"`
	TLS         string `md:"tls,required"`
	SchemaGen   string `md:"schemaGen,required"`
	Schema      string `md:"schema,required"`
}

func (s *Settings) ToMap() map[string]interface{} {

	properties := map[string]interface{}{
		"name":        s.Name,
		"description": s.Description,
		"version":     s.ApiVersion,
		"url":         s.URL,
		"tlsEnabled":  s.TLSEnabled,
		"user":        s.User,
		"password":    s.Password,
		"schemaGen":   s.SchemaGen,
	}

	if s.TLSEnabled {
		if "" != s.TLS {
			content, err := coerce.ToType(s.TLS, data.TypeObject)
			if nil == err {
				tlsBytes, err := b64.StdEncoding.DecodeString(strings.Split(content.(map[string]interface{})["content"].(string), ",")[1])
				if nil == err {
					properties["tls"] = string(tlsBytes)
				}
			}
		}
	}

	if "file" == s.SchemaGen {
		if "" != s.Schema {
			content, err := coerce.ToType(s.Schema, data.TypeObject)
			if nil == err {
				schema := content.(map[string]interface{})
				if nil != schema["content"] {
					schemaBytes, err := b64.StdEncoding.DecodeString(strings.Split(schema["content"].(string), ",")[1])
					if nil == err {
						properties["schema"] = string(schemaBytes)
					}
				}
			}
		}
	}

	return properties
}

func init() {
	if os.Getenv(log.EnvKeyLogLevel) == "DEBUG" {
		logCache.DebugEnabled()
	}

	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

// DgraphFactory for postgres connection
type DgraphFactory struct {
}

// Type DgraphFactory
func (this *DgraphFactory) Type() string {
	return "Dgraph"
}

// NewManager DgraphFactory
func (this *DgraphFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {

	s, err := NewSetting(settings)

	if err != nil {
		return nil, err
	}

	properties := s.ToMap()
	logCache.Debug("properties : ", properties)

	sharedConn := &SharedDgraphManager{
		name:       s.Name,
		properties: properties,
	}

	return sharedConn, nil
}

// SharedDgraphManager details
type SharedDgraphManager struct {
	name          string
	properties    map[string]interface{}
	dgraphService dbservice.UpsertService
}

func ReconstructGraph(graphData map[string]interface{}) model.Graph {
	return model.ReconstructGraph(graphData)
}

func (this *SharedDgraphManager) Lookup(clientID string, properties map[string]interface{}) (dbservice.UpsertService, error) {
	var err error
	if nil == this.dgraphService {
		for key, value := range this.properties {
			properties[key] = value
		}
		this.dgraphService, err = services.NewDgraphService(properties)
		logCache.Info("(DgraphServiceFactory.CreateUpsertService) upsertService : ", this.dgraphService)
		if nil != err {
			logCache.Info("(DgraphServiceFactory.CreateUpsertService) err : ", err)
			return nil, err
		}
	}
	return this.dgraphService.(dbservice.UpsertService), nil
}

// Type SharedDgraphManager details
func (this *SharedDgraphManager) Type() string {
	return "Dgraph"
}

// GetConnection SharedDgraphManager details
func (this *SharedDgraphManager) GetConnection() interface{} {
	return this.dgraphService
}

// ReleaseConnection SharedDgraphManager details
func (this *SharedDgraphManager) ReleaseConnection(connection interface{}) {

}

// Start SharedDgraphManager details
func (this *SharedDgraphManager) Start() error {
	return nil
}

// Stop SharedDgraphManager details
func (this *SharedDgraphManager) Stop() error {
	logCache.Debug("Cleaning up Graph")

	return nil
}
