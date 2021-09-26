package connection

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"database/sql"

	"github.com/project-flogo/core/activity"

	//pg provides native postres dependency
	_ "github.com/lib/pq"
)

type (
	//Connection datastructure for storing PostgreSQL connection details
	Connection struct {
		DatabaseURL string `json:"databaseURL"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		DbName      string `json:"databaseName"`
		SSLMode     string `json:"sslmode"`
		db          *sql.DB
	}

	//Connector is a representation of connector.json metadata for the postgres connection
	Connector struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Title       string `json:"title"`
		Version     string `json:"version"`
		Type        string `json:"type"`
		Ref         string `json:"ref"`
		Settings    []struct {
			Name  string      `json:"name"`
			Type  string      `json:"type"`
			Value interface{} `json:"value"`
		} `json:"settings"`
	}

	//Query structure for SQL Queries
	Query struct {
		TableName string            `json:"tableName"`
		Cols      []string          `json:"columns"`
		Filters   map[string]string `json:"filters"`
	}

	//record is the one row of the ResultSet retrieved from database after execution of SQL Query
	record map[string]interface{}

	//ResultSet is an aggregation of SQL Query data records fetched from the database
	ResultSet struct {
		Record []*record `json:"records"`
	}

	//Input is a representation of acitivity's input parametres
	Input struct {
		Parameters map[string]interface{}   `json:"parameters,omitempty"`
		Values     []map[string]interface{} `json:"values,omitempty"`
	}

	//QueryActivity provides Activity metadata for Flogo
	QueryActivity struct {
		metadata *activity.Metadata
	}
)

// NewConnection returns a deserialized conneciton object it does not establish a
// connection with the database. The client needs to call Login to establish a
// connection
/***func NewConnection(connector interface{}) (*Connection, error) {

	genericConn, err := generic.NewConnection(connector)
	if err != nil {
		return nil, errors.New("Failed to load SQLServer connection configuration")
	}
	conn := &Connection{}
	conn.Host, err = data.CoerceToString(genericConn.GetSetting("host"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for host failed: %s", err)
	}
	// log.Debugf("getconnection processed host: %s", conn.Host)
	conn.Port, err = data.CoerceToInteger(genericConn.GetSetting("port"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for port failed: %s", err)
	}
	// log.Debugf("getconnection processed port: %d", conn.Port)
	conn.User, err = data.CoerceToString(genericConn.GetSetting("user"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for user failed: %s", err)
	}
	//log.Debugf("getconnection processed user: %s", conn.User)
	conn.Password, err = data.CoerceToString(genericConn.GetSetting("password"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for password failed: %s", err)
	}
	log.Debugf("getconnection processed databaseName: %s", conn.DbName)
	conn.DbName, err = data.CoerceToString(genericConn.GetSetting("databaseName"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for databaseName failed: %s", err)
	}
	return conn, nil
}***/

//GetConnection returns a deserialized conneciton object it does not establish a
//connection with the database. The client needs to call Login to establish a
//connection
/***func GetConnection(connector interface{}) (*Connection, error) {

	genericConn, err := generic.NewConnection(connector)
	if err != nil {
		return nil, errors.New("Failed to load SQLServer connection configuration")
	}
	conn := &Connection{}
	conn.Host, err = data.CoerceToString(genericConn.GetSetting("host"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for host failed: %s", err)
	}
	log.Debugf("getconnection processed host: %s", conn.Host)
	conn.Port, err = data.CoerceToInteger(genericConn.GetSetting("port"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for port failed: %s", err)
	}
	log.Debugf("getconnection processed port: %d", conn.Port)
	conn.User, err = data.CoerceToString(genericConn.GetSetting("user"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for user failed: %s", err)
	}
	log.Debugf("getconnection processed user: %s", conn.User)
	conn.Password, err = data.CoerceToString(genericConn.GetSetting("password"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for password failed: %s", err)
	}
	log.Debugf("getconnection processed databaseName: %s", conn.DbName)
	conn.DbName, err = data.CoerceToString(genericConn.GetSetting("databaseName"))
	if err != nil {
		return nil, fmt.Errorf("connection getter for databaseName failed: %s", err)
	}

	err = conn.validate()
	if err != nil {
		return nil, fmt.Errorf("Connection validation error %s", err.Error())
	}

	return conn, nil
}***/

//validate validates and  set config values to connection struct. It returns an error
//if required value is not provided or cannot be correctly converted
/***func (con *Connection) validate() (err error) {

	if con.Host == "" {
		return fmt.Errorf("Required parameter Host missing, %s", err.Error())
	}

	if con.Port == 0 {
		return fmt.Errorf("Required parameter Port missing, %s", err.Error())
	}

	if con.User == "" {
		return fmt.Errorf("Required parameter User missing, %s", err.Error())
	}

	if con.DbName == "" {
		return fmt.Errorf("Required parameter DbName missing, %s", err.Error())
	}

	if con.Password == "" {
		return fmt.Errorf("Required parameter Password missing, %s", err.Error())
	}
	return nil
}***/

//QueryHelper is a simple query parser for extraction of values from an insert statement
type QueryHelper struct {
	sqlString   string
	values      []string
	first       string
	last        string
	valuesToken string
}

// NewQueryHelper creates a new instance of QueryHelper
func NewQueryHelper(sql string) *QueryHelper {
	qh := &QueryHelper{
		sqlString:   sql,
		valuesToken: "VALUES",
	}
	return qh
}

// Compose reconstitutes the query and returns it
func (qp *QueryHelper) Compose() string {
	return qp.first + qp.valuesToken + " " + strings.Join(qp.values, ", ") + " " + qp.last
}

// ComposeWithValues reconstitutes the query with external values
func (qp *QueryHelper) ComposeWithValues(values []string) string {
	return qp.first + " " + qp.valuesToken + " " + strings.Join(values, ", ") + " " + qp.last
}

// Decompose parses the SQL string to extract values from a SQL statement
func (qp *QueryHelper) Decompose() (values []string) {
	//sql := `INSERT INTO distributors (did, name) values (1, 'Cheese', 9.99), (2, 'Bread', 1.99), (3, 'Milk', 2.99) RETURNING (SELECT name FROM instructor)`
	// parts := strings.Split(qp.sqlString, "VALUES") //what if nested statement has a values too, not supporting that at the moment
	upperCaseQuery := strings.ToUpper(qp.sqlString)
	index := strings.Index(upperCaseQuery, "VALUES")
	if index == -1 {
		qp.first = qp.sqlString
		return
	}
	qp.first = qp.sqlString[0 : index-1]
	spart := qp.sqlString[index+len("values"):]
	spartLength := len(spart)

	i := 0
	braketCount := 0
	for i < spartLength {
		ch := spart[i]
		i++
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			continue
		}
		if ch == '(' {
			braketCount = braketCount + 1
			position := i
			for i < len(spart) {
				ch = spart[i]
				//fmt.Print(string(ch))
				i++
				if ch == '(' {
					braketCount = braketCount + 1
				}
				if ch == ')' || ch == 0 {
					braketCount = braketCount - 1
					if braketCount == 0 {
						break
					}
				}
			}
			value := "(" + spart[position:i-1] + ")"
			qp.values = append(qp.values, value)
			if i == spartLength {
				break
			}
			ch = spart[i]
			i++
			for ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
				ch = spart[i]
				i++
			}
			if ch != ',' {
				break
			}
		}
	}
	qp.last = spart[i-1:]
	return qp.values
}

func decodeBlob(blob string) []byte {
	decodedBlob, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return []byte(blob)
	}
	return decodedBlob
}

func (p *PgSharedConfigManager) getSchema(fields interface{}) map[string]string {
	schema := map[string]string{}
	for _, fieldObject := range fields.([]interface{}) {
		if fieldName, ok := fieldObject.(map[string]interface{})["FieldName"]; ok {
			schema[fieldName.(string)] = fieldObject.(map[string]interface{})["Type"].(string)
		}
	}
	return schema
}

//PreparedInsert allows querying database with named parameters
func (p *PgSharedConfigManager) PreparedInsert(queryString string, inputData *Input, fields interface{}) (results *ResultSet, err error) {

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

// UnmarshalRows function
func UnmarshalRows(rows *sql.Rows) (results *ResultSet, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column information, %s", err.Error())
	}

	count := len(columns)
	cols := make([]interface{}, count)
	args := make([]interface{}, count)
	coltypes, err := rows.ColumnTypes()

	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, fmt.Errorf("error determining column types, %s", err.Error())
	}

	for i := range cols {
		args[i] = &cols[i]
	}

	var resultSet ResultSet
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			fmt.Printf("%s", err.Error())
			return nil, fmt.Errorf("error scanning rows, %s", err.Error())
		}
		m := make(record)
		for i, b := range cols {
			dbType := coltypes[i].DatabaseTypeName()
			if b == nil {
				m[columns[i]] = nil
				continue
			}
			switch dbType {
			case "NUMERIC":
				x := b.([]uint8)
				if nx, ok := strconv.ParseFloat(string(x), 64); ok == nil {
					m[columns[i]] = nx
				}
			case "SMALLINT", "MEDIUMINT", "INT", "BIGINT":
				m[columns[i]] = b.(int64)
			case "DOUBLE", "REAL":
				m[columns[i]] = b.(float64)
			case "CHAR", "NCHAR", "BPCHAR", "":
				x := b.([]byte)
				if len, ok := coltypes[i].Length(); ok == true {
					if len > 0 {
						m[columns[i]] = string(x[0:len])
					}
				}
			case "BYTEA":
				m[columns[i]] = base64.StdEncoding.EncodeToString(b.([]byte))
			default:
				m[columns[i]] = b
			}
		}
		if len(m) > 0 {
			resultSet.Record = append(resultSet.Record, &m)
		}
	}
	return &resultSet, nil
}

//PreparedUpdate allows updating database with named parameters
func (p *PgSharedConfigManager) PreparedUpdate(queryString string, inputData *Input, fields interface{}) (results map[string]interface{}, err error) {
	logCache.Debugf("Executing prepared query %s", queryString)
	logCache.Debugf("Query parameters: %v", inputData)

	prepared := queryString
	schema := p.getSchema(fields)

	i := 0
	argCount := len(inputData.Parameters)
	inputArgs := make([]interface{}, argCount)
	for parameter, substitution := range inputData.Parameters {
		prepared = strings.Replace(prepared, "?"+string(parameter), "$"+strconv.Itoa(i+1), -1)
		//	log.Infof("adding parameter %d: %s", i, substitution)
		parameterType, ok := schema[parameter]
		if ok && parameterType == "BYTEA" {
			substitution = decodeBlob(substitution.(string))
		}
		inputArgs[i] = substitution
		i++
	}

	// log.Debugf("Prepared Query [%s]      Parameters [%v] ", prepared, inputArgs)
	stmt, err := p.db.Prepare(prepared)
	if err != nil {
		// log.Warnf("query preparation failed: %s, %s", prepared, err.Error())
		return nil, fmt.Errorf("query preparation failed: %s, %s", queryString, err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(inputArgs...)
	if err != nil {
		logCache.Errorf("PreparedQuery got error: %s", err)
		stmt.Close()
		return nil, err
	}

	output := make(map[string]interface{})
	logCache.Debugf("Number of rows affected: %d", result.RowsAffected)
	output["rowsAffected"], _ = result.RowsAffected()
	return output, nil
}

func checkCount(rows *sql.Rows) (count int, err error) {
	logCache.Debugf("Inside check count rows for update")
	// var counter int
	// defer rows.Close()
	for rows.Next() {
		logCache.Debugf("Row found")
		if err := rows.Scan(&count); err != nil {
			//log.Fatal(err)
		}
		logCache.Debugf("Counter: %d", count)
	}
	return count, err
}

//PreparedQuery allows querying database with named parameters
func (p *PgSharedConfigManager) PreparedQuery(queryString string, inputData *Input) (results *ResultSet, err error) {
	logCache.Infof("executing prepared query %s", queryString)

	prepared := queryString
	argCount := len(inputData.Parameters)
	inputArgs := make([]interface{}, argCount)

	i := 0
	for keyfield, value := range inputData.Parameters {
		prepared = strings.Replace(prepared, "?"+string(keyfield), "$"+strconv.Itoa(i+1), -1)
		logCache.Infof("adding parameter %d: %s", i, value)
		inputArgs[i] = value
		i++
	}
	logCache.Debugf("prepared query: %s ", prepared)
	stmt, err := p.db.Prepare(prepared)
	if err != nil {
		logCache.Warnf("query preparation failed: %s, %s", prepared, err.Error())
		return nil, fmt.Errorf("query preparation failed: %s, %s", queryString, err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(inputArgs...)

	if rows == nil {
		logCache.Infof("no rows returned for query %s", prepared)
		return nil, nil
	}

	defer rows.Close()
	logCache.Info("Return from PreparedQuery")
	return UnmarshalRows(rows)
}

// Login connects to the the postgres database cluster using the connection
// details provided in Connection configuration
func (con *Connection) Login() (err error) {
	if con.db != nil {
		return nil
	}

	conninfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		con.Host, con.Port, con.User, con.Password, con.DbName)

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return fmt.Errorf("Could not open connection to database %s, %s", con.DbName, err.Error())
	}
	con.db = db

	err = db.Ping()
	if err != nil {
		return err
	}

	// log.Infof("login successful")
	return nil
}

//Logout the database connection
func (con *Connection) Logout() (err error) {
	if con.db == nil {
		return nil
	}
	err = con.db.Close()
	// log.Infof("Logged out %s to %s", con.User, con.DbName)
	return
}
