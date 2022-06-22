/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package accumulator

import (
	"errors"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

const (
	lastElement = "LastElement"
)

type Settings struct {
	ArrayMode  bool `md:"ArrayMode"`
	WindowSize int  `md:"WindowSize"`
}

type Input struct {
	Input interface{} `md:"Input"`
}

type Output struct {
	Output []interface{} `md:"Output"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Input": i.Input,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	ok := true
	i.Input, ok = values["Input"]
	if !ok {
		return errors.New("Illegal Input type, expect map[string]interface{}.")
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Output": o.Output,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Output = values["Output"].([]interface{})
	return nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	_ = activity.Register(&Activity{}, New)
}

type Activity struct {
	arrayMode bool
	window    *Window
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	activity := &Activity{
		arrayMode: settings.ArrayMode,
		window:    NewWindow(settings.WindowSize),
	}
	return activity, nil
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	log := ctx.Logger()
	log.Info("[DataSelector:Eval] entering ........ ")
	defer log.Info("[DataSelector:Eval] exit ........ ")

	input := &Input{}
	ctx.GetInputObject(input)

	log.Info("Input data = ", input.Input)

	var outputTuple []interface{}
	if a.arrayMode {
		rawOutputTuple := input.Input.([]interface{})
		outputTuple = make([]interface{}, len(rawOutputTuple))
		for index, tuple := range rawOutputTuple {
			outputTuple[index] = tuple
			if index < len(rawOutputTuple)-1 {
				outputTuple[index].(map[string]interface{})[lastElement] = false
			} else {
				outputTuple[index].(map[string]interface{})[lastElement] = true
			}
		}
		log.Info("(Array Mode) Output data = ", outputTuple)
	} else {
		if nil != err {
			log.Error(err)
			return false, err
		}

		outputTuple, err = a.window.update(input.Input.(map[string]interface{}))
		if nil != err {
			log.Error(err)
			return false, err
		}
		log.Info("(Iterator Mode) Output data = ", outputTuple)
	}

	if nil != outputTuple {
		log.Debug("Output data = ", outputTuple)
		output := &Output{}
		output.Output = outputTuple
		err = ctx.SetOutputObject(output)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	return true, nil
}

func NewWindow(maxSize int) *Window {
	if 0 >= maxSize {
		maxSize = 1
	}
	window := &Window{
		currentSize: 0,
		maxSize:     maxSize,
		tuples:      make([]interface{}, maxSize),
	}
	return window
}

type Window struct {
	currentSize int
	maxSize     int
	tuples      []interface{}
}

func (this *Window) update(tuple map[string]interface{}) ([]interface{}, error) {
	this.currentSize += 1
	this.tuples[this.currentSize-1] = tuple
	if this.currentSize >= this.maxSize {
		tuple[lastElement] = true
		this.currentSize = 0
		return this.tuples, nil
	} else {
		tuple[lastElement] = false
	}
	return nil, nil
}
