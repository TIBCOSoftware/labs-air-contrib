package air

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/labs-air-contrib/common/util"
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
	template := params[3].(string)
	if !strings.Contains(template, "@") {
		return template, nil
	}

	dataMap := make(map[string]interface{})
	dataMap["f1..gateway"] = params[0]                           // gateway
	for key, value := range params[1].(map[string]interface{}) { // reading
		dataMap[fmt.Sprintf("f1..%s", key)] = value
	}
	if nil != params[2] { // enriched
		for _, element := range params[2].([]interface{}) {
			enrichedElement := element.(map[string]interface{})
			dataMap[fmt.Sprintf("%s..%s", enrichedElement["producer"], enrichedElement["name"])] = enrichedElement["value"]
		}
	}

	log.Debug("(fnAirDataSelector.Eval) input dataMap : ", dataMap)
	log.Debug("(fnAirDataSelector.Eval) input template : ", template)

	var data interface{}
	if '@' == template[0] && '@' == template[len(template)-1] {
		data = util.ExtractData(dataMap, template[1:len(template)-1])
	}
	if nil == data {
		data = NewKeywordMapper("@", "@").Replace(
			template,
			NewKeywordReplaceHandler(dataMap),
		)
	}

	log.Debug("(fnAirDataSelector.Eval) out data : ", data)

	if nil == data {
		return "", nil
	}
	return data, nil
}

func NewKeywordReplaceHandler(
	dataMap map[string]interface{},
) KeywordReplaceHandler {
	return KeywordReplaceHandler{
		dataMap: dataMap,
	}
}

type KeywordReplaceHandler struct {
	result  string
	dataMap map[string]interface{}
}

func (this *KeywordReplaceHandler) startToMap() {
	this.result = ""
}

func (this *KeywordReplaceHandler) Replace(keyword string) string {
	log.Debug("(KeywordReplaceHandler.Replace) keyword : ", keyword)
	log.Debug("(KeywordReplaceHandler.Replace) done .... ")
	return util.ExtractDataAsString(this.dataMap, keyword)
}

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
