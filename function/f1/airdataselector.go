package f1

import (
	"bytes"
	"fmt"
	"strings"

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
	reading := params[1].(map[string]interface{})
	reading["gateway"] = params[0]
	enriched := make(map[string]interface{})
	for _, element := range params[2].([]interface{}) {
		enrichedElement := element.(map[string]interface{})
		enriched[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
	}
	format := params[3].(string)

	log.Debug("(fnAirDataSelector.Eval) in reading : ", reading)
	log.Info("(fnAirDataSelector.Eval) in enriched : ", enriched)
	log.Info("(fnAirDataSelector.Eval) in format : ", format)

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
	keyElements := strings.Split(keyword, ".")
	if "f1" == keyElements[0] {
		data, ok := this.reading[keyElements[2]].(string)
		if ok {
			return strings.ReplaceAll(data, "\"", "\\\"")
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

//func (this *KeywordReplaceHandler) Replace(keyword string) string {
//	keyElements := strings.Split(keyword, ".")
//	if "f1" == keyElements[0] {
//		return fmt.Sprintf("%v", this.reading[keyElements[2]])
//	} else {
//		data := this.enriched[keyword]
//		if nil != data {
//			return fmt.Sprintf("%v", data)
//		}
//	}
//	return ""
//}

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
		//log.Infof("template[%d] = ", i, template[i])
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
				//log.Info("value ->", svalue);
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

		//log.Info("document = ", this.document)

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
