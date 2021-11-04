/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package builder

import (
	"fmt"

	"github.com/TIBCOSoftware/labs-air-contrib/connector/graph"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

const (
	PassThroughDataOut = "PassThroughDataOut"
)

var activityMd = activity.ToMetadata(&Input{}, &Output{})

func init() {

	_ = activity.Register(&Activity{}, New)
}

// Activity is the structure for Activity Metadata
type Activity struct {
	graphMgr       *graph.SharedGraphManager
	allowNullKey   bool
	batchMode      bool
	passThrough    []interface{}
	multiinstances string
}

// New for Activity
func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	sharedmanager := settings.GraphModel.(*graph.SharedGraphManager)
	act := &Activity{
		graphMgr:       sharedmanager,
		allowNullKey:   settings.AllowNullKey,
		batchMode:      settings.BatchMode,
		passThrough:    settings.PassThrough,
		multiinstances: settings.Multiinstances,
	}
	return act, nil
}

// Metadata  returns PostgreSQL's Query activity's meteadata
func (*Activity) Metadata() *activity.Metadata {
	return activityMd
}

var logCache = log.ChildLogger(log.RootLogger(), "labs-lc-activity-graphbuilder")

//Eval handles the data processing in the activity
func (a *Activity) Eval(context activity.Context) (done bool, err error) {

	logCache.Debug("[BuilderActivity:Eval] entering ........ ")
	logCache.Debug("[BuilderActivity:Eval] Exit ........ ")

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, err
	}

	logCache.Debug("[BuilderActivity:Eval] BatchEnd : ", input.BatchEnd)
	logCache.Info("[BuilderActivity:Eval] Nodes : ", input.Nodes)
	logCache.Info("[BuilderActivity:Eval] Edges : ", input.Edges)

	_, err = a.graphMgr.BuildGraph(input.Nodes, input.Edges, a.allowNullKey)

	if nil != err {
		return false, err
	}

	if false == a.batchMode || input.BatchEnd {

		graphData := a.graphMgr.ExportGraph()
		logCache.Debug("[BuilderActivity:Eval] Graph : ", graphData)
		context.SetOutput("Graph", graphData)

		passThroughDataDef := a.buildPassThroughData(context)
		if 0 != len(passThroughDataDef) {
			logCache.Info("[BuilderActivity:Eval] PassThroughData : ", input.PassThroughData)
			passThroughData := input.PassThroughData.(map[string]interface{})
			passThroughDataOut := make(map[string]interface{})
			for name, attrDef := range passThroughDataDef {
				value := passThroughData[name]
				defaultDV := attrDef.GetDValue()
				logCache.Debug("[BuilderActivity:Eval] name : ", name, ", value : ", value, ", default : ", defaultDV, ", optional : ", attrDef.IsOptional())
				if nil == value && !attrDef.IsOptional() {
					if nil != defaultDV {
						value = defaultDV
					} else {
						return false, fmt.Errorf("Data (%s)  should not be nil!", name)
					}

				}
				passThroughDataOut[name] = value
			}
			context.SetOutput(PassThroughDataOut, passThroughDataOut)
		}
		/* clear graph data */
	}

	return true, nil
}

func (a *Activity) buildPassThroughData(context activity.Context) map[string]*Field {
	passThroughData := make(map[string]*Field)
	logCache.Info("Processing handlers : PassThroughData = ", a.passThrough)

	for _, passThroughFieldname := range a.passThrough {
		passThroughFieldnameInfo := passThroughFieldname.(map[string]interface{})
		attribute := &Field{}
		attribute.SetName(passThroughFieldnameInfo["FieldName"].(string))
		attribute.SetType(passThroughFieldnameInfo["Type"].(string))
		attribute.SetOptional(nil != passThroughFieldnameInfo["Optional"] && "no" == passThroughFieldnameInfo["Optional"].(string))
		if nil != passThroughFieldnameInfo["Default"] && "" != passThroughFieldnameInfo["Default"].(string) {
			//attribute.SetDValue()
		}
		passThroughData[attribute.GetName()] = attribute
	}
	return passThroughData
}

type Field struct {
	name     string
	dValue   interface{}
	dataType string
	optional bool
}

func (this *Field) SetName(name string) {
	this.name = name
}

func (this *Field) GetName() string {
	return this.name
}

func (this *Field) SetDValue(dValue string) {
	this.dValue = dValue
}

func (this *Field) GetDValue() interface{} {
	return this.dValue
}

func (this *Field) SetType(dataType string) {
	this.dataType = dataType
}

func (this *Field) GetType() string {
	return this.dataType
}

func (this *Field) SetOptional(optional bool) {
	this.optional = optional
}

func (this *Field) IsOptional() bool {
	return this.optional
}
