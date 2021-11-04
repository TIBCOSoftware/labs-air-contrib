package dgraphupsert

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
	DgraphConnection   connection.Manager `md:"dgraphConnection,required"`
	CacheSize          int                `md:"cacheSize,required"`
	ReadableExternalId bool               `md:"readableExternalId,required"`
	ExplicitType       bool               `md:"explicitType,required"`
	TypeTag            string             `md:"typeTag,required"`
	AttrWithPrefix     bool               `md:"attrWithPrefix,required"`
}

func (s *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"cacheSize":          s.CacheSize,
		"readableExternalId": s.ReadableExternalId,
		"explicitType":       s.ExplicitType,
		"typeTag":            s.TypeTag,
		"addPrefixToAttr":    s.AttrWithPrefix,
	}
}

// Input Structure
type Input struct {
	Graph interface{} `md:"Graph,required"`
}

// ToMap Input interface
func (o *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Graph": o.Graph,
	}
}

// FromMap Input interface
func (o *Input) FromMap(values map[string]interface{}) error {
	var err error
	o.Graph, err = coerce.ToObject(values["Graph"])
	if err != nil {
		return err
	}

	return nil

}

//Output struct
type Output struct {
	MessageId interface{} `md:"MessageId,required"`
}

// ToMap conversion
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"MessageId": o.MessageId,
	}
}

// FromMap conversion
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	o.MessageId, err = coerce.ToString(values["MessageId"])
	if err != nil {
		return err
	}

	return nil
}
