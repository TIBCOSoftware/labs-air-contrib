package builder

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
	GraphModel     connection.Manager `md:"GraphModel,required"`
	AllowNullKey   bool               `md:"AllowNullKey,required"`
	BatchMode      bool               `md:"BatchMode,required"`
	PassThrough    []interface{}      `md:"PassThrough"`
	Multiinstances string             `md:"Multiinstances"`
}

// Input Structure
type Input struct {
	Nodes           interface{} `md:"Nodes,required"`
	Edges           interface{} `md:"Edges,required"`
	PassThroughData interface{} `md:"PassThroughData,required"`
	BatchEnd        bool        `md:"BatchEnd"`
}

// ToMap Input interface
func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Nodes":           o.Nodes,
		"Edges":           o.Edges,
		"PassThroughData": o.PassThroughData,
		"BatchEnd":        o.BatchEnd,
	}
}

// FromMap Input interface
func (o *Input) FromMap(values map[string]interface{}) error {
	var err error
	o.Nodes, err = coerce.ToObject(values["Nodes"])
	if err != nil {
		return err
	}

	o.Edges, err = coerce.ToObject(values["Edges"])
	if err != nil {
		return err
	}

	o.PassThroughData, err = coerce.ToObject(values["PassThroughData"])
	if err != nil {
		return err
	}

	o.BatchEnd, err = coerce.ToBool(values["BatchEnd"])
	if err != nil {
		return err
	}

	return nil

}

//Output struct
type Output struct {
	Graph              interface{} `md:"Graph,required"`
	PassThroughDataOut interface{} `md:"PassThroughDataOut,required"`
}

// ToMap conversion
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Graph":              o.Graph,
		"PassThroughDataOut": o.PassThroughDataOut,
	}
}

// FromMap conversion
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	o.Graph, err = coerce.ToObject(values["Graph"])
	if err != nil {
		return err
	}

	o.PassThroughDataOut, err = coerce.ToObject(values["PassThroughDataOut"])
	if err != nil {
		return err
	}

	return nil
}
