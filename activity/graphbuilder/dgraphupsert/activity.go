/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package dgraphupsert

import (
	"github.com/TIBCOSoftware/labs-air-contrib/connector/dgraph"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

const (
	Connection         = "dgraphConnection"
	cacheSize          = "cacheSize"
	typeTag            = "typeTag"
	explicitType       = "explicitType"
	readableExternalId = "readableExternalId"
	attrWithPrefix     = "attrWithPrefix"
)

var activityMd = activity.ToMetadata(&Input{}, &Output{})

func init() {

	_ = activity.Register(&Activity{}, New)
}

// Activity is the structure for Activity Metadata
type Activity struct {
	dgraphMgr          *dgraph.SharedDgraphManager
	settings           map[string]interface{}
	cacheSize          int
	readableExternalId bool
	explicitType       bool
	typeTag            string
	attrWithPrefix     bool
}

// New for Activity
func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	sharedmanager := settings.DgraphConnection.(*dgraph.SharedDgraphManager)
	act := &Activity{
		dgraphMgr:          sharedmanager,
		settings:           settings.ToMap(),
		cacheSize:          settings.CacheSize,
		readableExternalId: settings.ReadableExternalId,
		explicitType:       settings.ExplicitType,
		typeTag:            settings.TypeTag,
		attrWithPrefix:     settings.AttrWithPrefix,
	}
	return act, nil
}

// Metadata  returns PostgreSQL's Query activity's meteadata
func (*Activity) Metadata() *activity.Metadata {
	return activityMd
}

var logCache = log.ChildLogger(log.RootLogger(), "labs-lc-activity-dgraphupsert")

//Eval handles the data processing in the activity
func (a *Activity) Eval(context activity.Context) (done bool, err error) {

	logCache.Debug("[DgraphUpsertActivity:Eval] entering ........ ")
	logCache.Debug("[DgraphUpsertActivity:Eval] Exit ........ ")

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, err
	}

	logCache.Debug("[BuilderActivity:Eval] Graph : ", input.Graph)

	graph := dgraph.ReconstructGraph(input.Graph.(map[string]interface{})["graph"].(map[string]interface{}))

	logCache.Debug("(getDgraph) graph obj = ", graph)

	a.settings["graphModel"] = graph.GetModel()
	dgraphService, err := a.dgraphMgr.Lookup(context.Name(), a.settings)
	err = dgraphService.UpsertGraph(graph, nil)

	if nil != err {
		logCache.Error("(DgraphUpsertActivity) exit during upsert, with error = ", err.Error())
		return false, err
	}

	return true, nil
}
