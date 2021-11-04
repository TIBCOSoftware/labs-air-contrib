package log

import (
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
)

func init() {
	_ = activity.Register(&Activity{})
}

type Input struct {
	LogLevel   string `md:"Log Level"`  //
	Message    string `md:"message"`    // The message to log
	FlowInfo   bool   `md:"flowInfo"`   //
	AddDetails bool   `md:"addDetails"` // Append contextual execution information to the log message
	UsePrint   bool   `md:"usePrint"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Log Level":  i.LogLevel,
		"message":    i.Message,
		"flowInfo":   i.FlowInfo,
		"addDetails": i.AddDetails,
		"usePrint":   i.UsePrint,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.LogLevel, err = coerce.ToString(values["Log Level"])
	if err != nil {
		return err
	}
	i.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}
	i.FlowInfo, err = coerce.ToBool(values["flowInfo"])
	if err != nil {
		return err
	}
	i.AddDetails, err = coerce.ToBool(values["addDetails"])
	if err != nil {
		return err
	}

	i.UsePrint, err = coerce.ToBool(values["usePrint"])
	if err != nil {
		return err
	}

	return nil
}

var activityMd = activity.ToMetadata(&Input{})

// Activity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	ctx.GetInputObject(input)

	msg := input.Message

	if input.AddDetails {
		msg = fmt.Sprintf("'%s' - HostID [%s], HostName [%s], Activity [%s]", msg,
			ctx.ActivityHost().ID(), ctx.ActivityHost().Name(), ctx.Name())
	}

	if input.UsePrint {
		fmt.Println(msg)
	} else {
		switch input.LogLevel {
		case "INFO":
			ctx.Logger().Info(msg)
		case "DEBUG":
			ctx.Logger().Debug(msg)
		}
	}

	return true, nil
}
