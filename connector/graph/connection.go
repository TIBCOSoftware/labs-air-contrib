package graph

import (
	"errors"
	"os"

	"github.com/SteveNY-Tibco/labs-air-contrib/common/graphbuilder/model"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "Postgres.connection")
var factory = &GraphFactory{graphBuilder: model.NewGraphBuilder()}

// Settings for postgres
type Settings struct {
	Name        string `md:"name,required"`
	Description string `md:"description"`
	ModelSource string `md:"modelSource,required"`
	URL         string `md:"url"`
	Model       string `md:"model,required"`
	Metadata    string `md:"metadata,required"`
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

// GraphFactory for postgres connection
type GraphFactory struct {
	graphBuilder *model.GraphBuilder
}

// Type GraphFactory
func (this *GraphFactory) Type() string {
	return "Graph"
}

// NewManager GraphFactory
func (this *GraphFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {

	sharedConn := &SharedGraphManager{}
	var err error

	s := &Settings{}

	err = metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}

	cName := s.Name
	if cName == "" {
		return nil, errors.New("Required Parameter Name is missing")
	}

	cModelSource := s.ModelSource
	if cModelSource == "" {
		return nil, errors.New("Required Parameter ModelSource is missing")
	}

	cModel := s.Model
	if cModel == "" {
		return nil, errors.New("Required Parameter Model is missing")
	}

	cMetadata := s.Metadata
	if cMetadata == "" {
		return nil, errors.New("Required Parameter Metadata is missing")
	}

	model, err := model.NewGraphModel(cName, cMetadata)
	if nil != err {
		return nil, err
	}

	sharedConn.name = cName
	sharedConn.model = model

	return sharedConn, nil
}

// SharedGraphManager details
type SharedGraphManager struct {
	name  string
	model *model.GraphDefinition
}

// Create graph
func (this *SharedGraphManager) CreateGraph(
	nodes interface{},
	edges interface{},
	allowNullKey bool) (*model.Graph, error) {

	graphId := this.model.GetId()
	deltaGraph := factory.graphBuilder.CreateGraph(graphId, this.model)
	err := factory.graphBuilder.BuildGraph(
		&deltaGraph,
		this.model,
		nodes,
		edges,
		allowNullKey,
	)

	if nil != err {
		return nil, err
	}

	return &deltaGraph, nil
}

// Type SharedGraphManager details
func (this *SharedGraphManager) Type() string {
	return "Graph"
}

// GetConnection SharedGraphManager details
func (this *SharedGraphManager) GetConnection() interface{} {
	return this
}

// ReleaseConnection SharedGraphManager details
func (this *SharedGraphManager) ReleaseConnection(connection interface{}) {

}

// Start SharedGraphManager details
func (this *SharedGraphManager) Start() error {
	return nil
}

// Stop SharedGraphManager details
func (this *SharedGraphManager) Stop() error {
	logCache.Debug("Cleaning up Graph")

	return nil
}
