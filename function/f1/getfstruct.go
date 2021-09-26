package f1

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnGetFStruct{})
}

type fnGetFStruct struct {
}

func (fnGetFStruct) Name() string {
	return "getfstruct"
}

func (fnGetFStruct) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeInt}, false
}

func (fnGetFStruct) Eval(params ...interface{}) (interface{}, error) {
	folder, ok1 := params[0].(string)
	if !ok1 {
		return nil, fmt.Errorf("Illegal parameter : folder string")
	}

	maxDepth, ok2 := params[1].(int)
	if !ok2 {
		return nil, fmt.Errorf("Illegal parameter : depth int")
	}

	result := walk(folder, 0, maxDepth)

	return result, nil
}

func walk(filename string, currentDepth int, maxDepth int) []interface{} {
	//fmt.Println("name ---> ", filename)
	result := make([]interface{}, 0)
	if currentDepth <= maxDepth {
		if stat, err := os.Stat(filename); err == nil {
			switch mode := stat.Mode(); {
			case mode.IsDir():
				//fmt.Println("directory")
				files, _ := ioutil.ReadDir(filename)
				for index := range files {
					result = append(result, map[string]interface{}{
						"Name":  files[index].Name(),
						"Type":  "folder",
						"Value": walk(fmt.Sprintf("%s/%s", filename, files[index].Name()), currentDepth+1, maxDepth),
					})
				}

			case mode.IsRegular():
				result = append(result, map[string]interface{}{
					"Name":  filename,
					"Type":  "file",
					"Value": "",
				})
			}
		} else if os.IsNotExist(err) {

		}
	}
	return result
}
