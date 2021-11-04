package notificationlistener

type Settings struct {
}

type HandlerSettings struct {
	NotifierID string `md:"notifierID,required"`
}

type Output struct {
	Gateway  string                 `md:"gateway,required"`
	Reading  map[string]interface{} `md:"reading"`
	Enriched []interface{}          `md:"enriched"`
}

func (this *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"gateway":  this.Gateway,
		"reading":  this.Reading,
		"enriched": this.Enriched,
	}
}

func (this *Output) FromMap(values map[string]interface{}) error {

	this.Gateway = values["gateway"].(string)
	this.Reading = values["reading"].(map[string]interface{})
	this.Enriched = values["enriched"].([]interface{})

	return nil
}
