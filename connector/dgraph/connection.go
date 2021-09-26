package connection

import (
	b64 "encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/SteveNY-Tibco/labs-air-contrib/common/graphbuilder/dbservice"
	dbsf "github.com/SteveNY-Tibco/labs-air-contrib/common/graphbuilder/dbservice/factory"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "Dgraph.connection")
var factory = &DgraphFactory{}

// Settings for postgres
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

	sharedConn := &SharedDgraphManager{}
	s := &Settings{}

	var err = metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}

	cName := s.Name
	if cName == "" {
		return nil, errors.New("Required Parameter Name is missing")
	}

	cDescription := s.Description
	if cDescription == "" {
		return nil, errors.New("Required Parameter ModelSource is missing")
	}

	cApiVersion := s.ApiVersion
	if cApiVersion == "" {
		return nil, errors.New("Required Parameter Model is missing")
	}

	cURL := s.URL
	if cURL == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}

	cTLSEnabled := s.TLSEnabled

	cUser := s.User
	if cUser == "" {
		return nil, errors.New("Required Parameter Name is missing")
	}

	cPassword := s.Password
	if cPassword == "" {
		return nil, errors.New("Required Parameter ModelSource is missing")
	}

	cTLS := s.TLS
	if cTLS == "" {
		return nil, errors.New("Required Parameter Model is missing")
	}

	cSchemaGen := s.SchemaGen
	if cSchemaGen == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}

	cSchema := s.Schema
	if cSchema == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}

	properties := make(map[string]interface{})

	properties["version"] = cApiVersion
	properties["url"] = cURL
	properties["user"] = cUser
	properties["password"] = cPassword
	properties["tlsEnabled"] = cTLSEnabled
	if cTLSEnabled {
		if "" != cTLS {
			content, err := coerce.ToType(cTLS, data.TypeObject)
			if nil == err {
				tlsBytes, err := b64.StdEncoding.DecodeString(strings.Split(content.(map[string]interface{})["content"].(string), ",")[1])
				if nil == err {
					properties["tls"] = string(tlsBytes)
				}
			}
		}
	}

	properties["schemaGen"] = cSchemaGen
	if "file" == cSchemaGen {
		if "" != cSchema {
			content, err := coerce.ToType(cSchema, data.TypeObject)
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

	logCache.Debug("properties : ", properties)

	_, err = dbsf.GetFactory(dbservice.Dgraph).CreateUpsertService(cName, properties)
	//dgraph.GetFactory().CreateService(connectorName, properties)
	//sharedConn.dgraphService = dgraphService

	if nil != err {
		return nil, err
	}

	return sharedConn, nil
}

// SharedDgraphManager details
type SharedDgraphManager struct {
	name          string
	properties    map[string]interface{}
	dgraphService dbservice.DBServiceFactory
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
