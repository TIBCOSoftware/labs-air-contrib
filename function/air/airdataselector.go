package air

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/TIBCOSoftware/labs-lightcrane-contrib/common/objectbuilder"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnAirDataSelector{})
}

type fnAirDataSelector struct {
}

func (fnAirDataSelector) Name() string {
	return "airdataselector"
}

func (fnAirDataSelector) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeObject, data.TypeArray, data.TypeString}, false
}

func (fnAirDataSelector) Eval(params ...interface{}) (interface{}, error) {
	// f1.airdataselector($flow.gateway, $flow.reading, $flow.enriched, $property["Python.DataIn"])
	reading := params[1].(map[string]interface{})
	reading["gateway"] = params[0]
	enriched := make(map[string]interface{})
	for _, element := range params[2].([]interface{}) {
		enrichedElement := element.(map[string]interface{})
		enriched[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
	}
	format := params[3].(string)

	log.Debug("(fnAirDataSelector.Eval) in reading : ", reading)
	log.Debug("(fnAirDataSelector.Eval) in enriched : ", enriched)
	log.Debug("(fnAirDataSelector.Eval) in format : ", format)

	data := NewKeywordMapper("@", "@").Replace(
		format,
		NewKeywordReplaceHandler(reading, enriched),
	)

	log.Debug("(fnAirDataSelector.Eval) out data string : ", data)

	return data, nil
}

func NewKeywordReplaceHandler(
	reading map[string]interface{},
	enriched map[string]interface{},
) KeywordReplaceHandler {
	return KeywordReplaceHandler{
		reading:  reading,
		enriched: enriched,
	}
}

type KeywordReplaceHandler struct {
	result   string
	reading  map[string]interface{}
	enriched map[string]interface{}
}

func (this *KeywordReplaceHandler) startToMap() {
	this.result = ""
}

func (this *KeywordReplaceHandler) Replace(keyword string) string {
	log.Debug("(KeywordReplaceHandler.Replace) keyword : ", keyword)
	keyElements := strings.Split(keyword, ".")
	subkeyElements := strings.Split(keyElements[2], "/")
	log.Debug("(KeywordReplaceHandler.Replace) real keyword : ", keyElements[2])
	var data interface{}
	if "f1" == keyElements[0] {
		data = this.reading[subkeyElements[0]]
	} else {
		/*
			keyword : PythonService1..Result/result[]
			keyElements[0] : PythonService1
			subkeyElements[0] : Result
			enriched : map[PythonService1..Result:{"id": "process:abc", "input1": [[2, 1], [3, 4]], "input2": [[6, 5], [8, 7]], "result": [2, 1, 3, 4, 6, 5, 8, 7]}]

		*/
		data = this.enriched[fmt.Sprintf("%s..%s", keyElements[0], subkeyElements[0])]
	}
	log.Debug("(KeywordReplaceHandler.Replace) real data : ", data)

	dataType := reflect.ValueOf(data).Kind()
	log.Debug("(KeywordReplaceHandler.Replace) dataType : ", dataType.String())
	if reflect.String == dataType {
		return strings.ReplaceAll(data.(string), "\"", "\\\"")
	} else if reflect.Map == dataType {
		if len(subkeyElements) > 1 {
			log.Debug("(KeywordReplaceHandler.Replace) keyElements[2] : ", keyElements[2])
			subkey := fmt.Sprintf("root%s", strings.Replace(keyElements[2][len(subkeyElements[0]):], "/", ".", -1))
			log.Debug("(KeywordReplaceHandler.Replace) subkey : ", subkey)
			data = objectbuilder.LocateObject(data.(map[string]interface{}), subkey).(interface{})
			log.Debug("(KeywordReplaceHandler.Replace) data : ", data)
		}
		realDataType := reflect.ValueOf(data).Kind()
		log.Debug("(KeywordReplaceHandler.Replace) realDataType : ", realDataType.String())
		if reflect.Map == realDataType || reflect.Array == realDataType || reflect.Slice == realDataType {
			jsonBuf, _ := json.Marshal(data)
			log.Debug("(KeywordReplaceHandler.Replace) string(jsonBuf) : ", string(jsonBuf))
			return fmt.Sprintf("%v", string(jsonBuf))
		} else {
			log.Debug("(KeywordReplaceHandler.Replace) data.(string) : ", data.(string))
			return strings.ReplaceAll(data.(string), "\"", "\\\"")
		}
	} else if reflect.Array == dataType {
		jsonBuf, _ := json.Marshal(data)
		return fmt.Sprintf("%v", string(jsonBuf))
	}
	return fmt.Sprintf("%v", data)
}

/*
func (this *KeywordReplaceHandler) Replace(keyword string) string {
	log.Debug("(KeywordReplaceHandler.Replace) keyword : ", keyword)
	keyElements := strings.Split(keyword, ".")
	if "f1" == keyElements[0] {
		subkeyElements := strings.Split(keyElements[2], "/")
		log.Debug("(KeywordReplaceHandler.Replace) real keyword : ", keyElements[2])
		data := this.reading[subkeyElements[0]]
		dataType := reflect.ValueOf(data).Kind()
		if reflect.String == dataType {
			return strings.ReplaceAll(data.(string), "\"", "\\\"")
		} else if reflect.Map == dataType {
			if len(subkeyElements) > 1 {
				log.Debug("(KeywordReplaceHandler.Replace) keyElements[2] : ", keyElements[2])
				subkey := fmt.Sprintf("root%s", strings.Replace(keyElements[2][len(subkeyElements[0]):], "/", ".", -1))
				log.Debug("(KeywordReplaceHandler.Replace) subkey : ", subkey)
				data = objectbuilder.LocateObject(data.(map[string]interface{}), subkey).(interface{})
				log.Debug("(KeywordReplaceHandler.Replace) data : ", data)
			}
			realDataType := reflect.ValueOf(data).Kind()
			log.Debug("(KeywordReplaceHandler.Replace) realDataType : ", realDataType.String())
			if reflect.Map == realDataType || reflect.Array == realDataType || reflect.Slice == realDataType {
				jsonBuf, _ := json.Marshal(data)
				log.Debug("(KeywordReplaceHandler.Replace) string(jsonBuf) : ", string(jsonBuf))
				return fmt.Sprintf("%v", string(jsonBuf))
			} else {
				log.Debug("(KeywordReplaceHandler.Replace) data.(string) : ", data.(string))
				return strings.ReplaceAll(data.(string), "\"", "\\\"")
			}
		} else if reflect.Array == dataType {
			jsonBuf, _ := json.Marshal(data)
			return fmt.Sprintf("%v", string(jsonBuf))
		}
		return fmt.Sprintf("%v", this.reading[keyElements[2]])
	} else {
		data := this.enriched[keyword]
		if nil != data {
			return fmt.Sprintf("%v", data)
		}
	}
	return ""
}

func (this *KeywordReplaceHandler) Replace(keyword string) string {
	keyElements := strings.Split(keyword, ".")
	if "f1" == keyElements[0] {
		data := this.reading[keyElements[2]]
		dataType := reflect.ValueOf(data).Kind()
		if reflect.String == dataType {
			return strings.ReplaceAll(data.(string), "\"", "\\\"")
		} else if reflect.Map == dataType || reflect.Array == dataType {
			jsonBuf, _ := json.Marshal(data)
			return fmt.Sprintf("%v", string(jsonBuf))
		}
		return fmt.Sprintf("%v", this.reading[keyElements[2]])
	} else {
		data := this.enriched[keyword]
		if nil != data {
			return fmt.Sprintf("%v", data)
		}
	}
	return ""
}

func (this *KeywordReplaceHandler) Replace(keyword string) string {
	keyElements := strings.Split(keyword, ".")
	if "f1" == keyElements[0] {
		return fmt.Sprintf("%v", this.reading[keyElements[2]])
	} else {
		data := this.enriched[keyword]
		if nil != data {
			return fmt.Sprintf("%v", data)
		}
	}
	return ""
}*/

func (this *KeywordReplaceHandler) endOfMapping(document string) {
	this.result = document
}

func (this *KeywordReplaceHandler) getResult() string {
	return this.result
}

func NewKeywordMapper(
	lefttag string,
	righttag string) *KeywordMapper {
	mapper := KeywordMapper{
		keywordOnly:  false,
		slefttag:     lefttag,
		srighttag:    righttag,
		slefttaglen:  len(lefttag),
		srighttaglen: len(righttag),
	}
	return &mapper
}

type KeywordMapper struct {
	keywordOnly  bool
	slefttag     string
	srighttag    string
	slefttaglen  int
	srighttaglen int
	document     bytes.Buffer
	keyword      bytes.Buffer
}

func (this *KeywordMapper) Replace(template string, mh KeywordReplaceHandler) string {
	if "" == template {
		return ""
	}

	log.Debug("[KeywordMapper.replace] template = ", template)

	this.document.Reset()
	this.keyword.Reset()

	scope := false
	boundary := false
	skeyword := ""
	svalue := ""

	mh.startToMap()
	for i := 0; i < len(template); i++ {
		//log.Debugf("template[%d] = ", i, template[i])
		// maybe find a keyword beginning Tag - now isn't in a keyword
		if !scope && template[i] == this.slefttag[0] {
			if this.isATag(i, this.slefttag, template) {
				this.keyword.Reset()
				scope = true
			}
		} else if scope && template[i] == this.srighttag[0] {
			// maybe find a keyword ending Tag - now in a keyword
			if this.isATag(i, this.srighttag, template) {
				i = i + this.srighttaglen - 1
				skeyword = this.keyword.String()[this.slefttaglen:this.keyword.Len()]
				svalue = mh.Replace(skeyword)
				if "" == svalue {
					svalue = fmt.Sprintf("%s%s%s", this.slefttag, skeyword, this.srighttag)
				}
				//log.Debug("value ->", svalue);
				this.document.WriteString(svalue)
				boundary = true
				scope = false
			}
		}

		if !boundary {
			if !scope && !this.keywordOnly {
				this.document.WriteByte(template[i])
			} else {
				this.keyword.WriteByte(template[i])
			}
		} else {
			boundary = false
		}

		//log.Debug("document = ", this.document)

	}
	mh.endOfMapping(this.document.String())
	return mh.getResult()
}

func (this *KeywordMapper) isATag(i int, tag string, template string) bool {
	if len(template) >= len(tag) {
		for j := 0; j < len(tag); j++ {
			if tag[j] != template[i+j] {
				return false
			}
		}
		return true
	}
	return false
}
