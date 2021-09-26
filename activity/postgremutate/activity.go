package postgremutate

import (
	"encoding/json"
	"fmt"

	"git.tibco.com/git/product/ipaas/wi-contrib.git/engine/jsonschema"
	"git.tibco.com/git/product/ipaas/wi-postgres.git/src/app/PostgreSQL/connector/connection"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"
)

/*** const (
	 ConnectionProp    = "Connection"
	 DatabaseURL       = "databaseURL"
	 Host              = "host"
	 Port              = "port"
	 User              = "username"
	 Password          = "password"
	 DatabaseName      = "databaseName"
	 InputProp         = "input"
	 ActivityOutput    = "output"
	 QueryProperty     = "Query"
	 RuntimeQuery      = "RuntimeQuery"
	 QueryNameProperty = "QueryName"
	 FieldsProperty    = "Fields"
	 OutputProperty    = "Output"
	 RecordsProperty   = "records"
 )***/

var activityMd = activity.ToMetadata(&Input{}, &Output{})

func init() {

	_ = activity.Register(&Activity{}, New)
}

// Activity is the structure for Activity Metadata
type Activity struct {
}

// New for Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

// Metadata  returns PostgreSQL's Query activity's meteadata
func (*Activity) Metadata() *activity.Metadata {
	return activityMd
}

var logCache = log.ChildLogger(log.RootLogger(), "postgres-insert")

//Eval handles the data processing in the activity
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	logCache.Info("Executing PostgreSQL Insert Activity")

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, err
	}
	sharedmanager := input.Connection.(*connection.PgSharedConfigManager)

	query := input.Query
	if query == "" {
		return false, fmt.Errorf("missing schema SQL statement")
	}

	inputParams, err := getInputData(input.InputParams)
	if err != nil {
		return false, fmt.Errorf("failed to read input arguments: %s", err.Error())
	}

	result, err := insert(query, inputParams, input.Fields)
	if err != nil {
		return false, fmt.Errorf("query execution failed: %s", err.Error())
	}
	//json schema validation
	outputSchema := input.Fields
	if outputSchema != nil {
		jsonBytes, err := json.Marshal(outputSchema)
		schema := string(jsonBytes)
		err = jsonschema.ValidateFromObject(schema, result)
		if err != nil {
			return false, fmt.Errorf("Schema validation error %s", err.Error())
		}
	}
	output := &Output{}
	output.Output = result
	err = context.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getInputData(inputData interface{}) (inputParams *connection.Input, err error) {

	inputParams = &connection.Input{}

	if inputData == nil {
		return nil, fmt.Errorf("missing input arguments")
	}

	switch inputData.(type) {
	case string:
		logCache.Debugf("Input data content: %s", inputData.(string))
		tempMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(inputData.(string)), &tempMap)
		if err != nil {
			return nil, fmt.Errorf("string parameter read error: %s", err.Error())
		}
		inputParams.Parameters = tempMap
	default:
		dataBytes, err := json.Marshal(inputData)
		logCache.Debugf("input arguments data: %s", string(dataBytes))
		if err != nil {
			return nil, fmt.Errorf("input data read failed: %s", err.Error())
		}
		err = json.Unmarshal(dataBytes, inputParams)
		if err != nil {
			return nil, fmt.Errorf("complex parameters read error:, %s", err.Error())
		}
	}
	return
}

//PreparedInsert allows querying database with named parameters
func insert(queryString string, inputData *Input, fields interface{}) (results *ResultSet, err error) {

	logCache.Debugf("executing prepared query %s", queryString)
	logCache.Debugf("inputParms: %v", inputData)

	schema := p.getSchema(fields)
	prepared := strings.Trim(queryString, " ")
	if queryString[len(queryString)-1] != ';' {
		prepared = prepared + ";"
	}
	queryHelper := NewQueryHelper(queryString)
	queryValues := queryHelper.Decompose()

	queryValuesLength := len(queryValues)
	inputValuesLength := len(inputData.Values)
	if queryValuesLength > 1 && queryValuesLength != inputValuesLength && inputValuesLength != 0 {
		return nil, fmt.Errorf("input data values length does not match query input data values length, %v != %v", queryValuesLength, inputValuesLength)
	}

	inputArgs := []interface{}{} // not very efficient better to size, but to find size you have to search, will optimize later
	replacedValues := []string{}
	argIndex := 1

	// Substitute the values
	if queryValuesLength > 0 {
		queryValue := queryValues[0]
		for inputValuesIndex, inputValues := range inputData.Values {
			if queryValuesLength > 1 {
				queryValue = queryValues[inputValuesIndex]
			}
			regExp := regexp.MustCompile("\\?\\w*")
			matches := regExp.FindAllStringSubmatch(queryValue, -1)
			value := queryValue
			for _, match := range matches {
				parameter := strings.Split(match[0], "?")[1]
				substitution, ok := inputData.Parameters[parameter]
				if !ok {
					substitution, ok = inputValues[parameter]
					if !ok {
						return nil, fmt.Errorf("missing substitution for: %s", match[0])
					}
				}
				// replace the first occurance, as it is found
				value = strings.Replace(value, match[0], "$"+strconv.Itoa(argIndex), 1)
				argIndex++
				// log.Debugf("prepared statement: %s", value)

				parameterType, ok := schema[parameter]
				if ok && parameterType == "BYTEA" {
					substitution = decodeBlob(substitution.(string))
				}
				inputArgs = append(inputArgs, substitution)
			}
			replacedValues = append(replacedValues, value)
		}
	}
	if len(replacedValues) > 0 {
		prepared = queryHelper.ComposeWithValues(replacedValues)
	}
	// Now substitute the parameters if there are any. Parameters always come as object keys
	// of the inputData.Parameters. We use the prepared statement now as all values, if any
	// have already been substituted and equivalent input arguments created
	regExp := regexp.MustCompile("\\?\\w*")
	matches := regExp.FindAllStringSubmatch(prepared, -1)
	for _, match := range matches {
		parameter := strings.Split(match[0], "?")[1]
		substitution := inputData.Parameters[parameter]
		if substitution == nil {
			return nil, fmt.Errorf("missing parameter substitution for: %s", match[0])
		}
		// replace the first occurance, as it is found
		prepared = strings.Replace(prepared, match[0], "$"+strconv.Itoa(argIndex), 1)
		argIndex++
		logCache.Debugf("prepared statement: %s", prepared)
		parameterType, ok := schema[parameter]
		if ok && parameterType == "BYTEA" {
			substitution = decodeBlob(substitution.(string))
		}
		inputArgs = append(inputArgs, substitution)
	}

	// log.Debug("prepared insert [%s]", prepared)
	stmt, err := p.db.Prepare(prepared)
	if err != nil {
		// log.Warnf("query preparation failed: %s, %s", prepared, err.Error())
		return nil, fmt.Errorf("query preparation failed: %s, %s", queryString, err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(inputArgs...)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		logCache.Infof("no rows returned for query %s", prepared)
		return nil, nil
	}
	defer rows.Close()
	logCache.Info("return from prepared insert")
	return UnmarshalRows(rows)
}
