package sqlinsert

import (
	"database/sql"
	"fmt"

	"github.com/TIBCOSoftware/labs-air-contrib/activity/sqlinsert/util"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

const (
	ovResults = "results"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{MaxIdleConns: 2}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	dbHelper, err := util.GetDbHelper(s.DbType)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("DB: '%s'", s.DbType)

	log.RootLogger().Debugf("New activity settings:  %v", s)

	// todo move this to a shared connection object
	db, err := getConnection(s)
	if err != nil {
		return nil, err
	}

	sqlStatement, err := util.NewSQLStatement(dbHelper, s.Statement)
	if err != nil {
		return nil, err
	}

	if sqlStatement.Type() != util.StInsert {
		return nil, fmt.Errorf("only insert statement is supported")
	}

	act := &Activity{db: db, dbHelper: dbHelper, sqlStatement: sqlStatement}

	if !s.DisablePrepared {
		ctx.Logger().Debugf("Using PreparedStatement: %s", sqlStatement.PreparedStatementSQL())
		act.stmt, err = db.Prepare(sqlStatement.PreparedStatementSQL())
		if err != nil {
			return nil, err
		}
	}

	return act, nil
}

// Activity is a Counter Activity implementation
type Activity struct {
	dbHelper       util.DbHelper
	db             *sql.DB
	sqlStatement   *util.SQLStatement
	stmt           *sql.Stmt
	labeledResults bool
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Cleanup() error {
	if a.stmt != nil {
		err := a.stmt.Close()
		log.RootLogger().Warnf("error cleaning up SQL Insert activity: %v", err)
	}

	log.RootLogger().Tracef("cleaning up SQL Insert activity")

	return a.db.Close()
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	in := &Input{}
	err = ctx.GetInputObject(in)
	if err != nil {
		return false, err
	}

	log.RootLogger().Debugf("Eval input params: %v  Total num params: %v", in.Params, len(in.Params))

	results, err := a.doInsert(in.Params)
	if err != nil {
		return false, err
	}

	err = ctx.SetOutput(ovResults, results)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *Activity) doInsert(params map[string]interface{}) (interface{}, error) {

	var err error
	var results sql.Result

	if a.stmt != nil {
		// MAG
		log.RootLogger().Debugf("Executing statement: %v", " state")
		// args := a.sqlStatement.GetPreparedStatementArgs(params)
		// rows, err = a.stmt.Exec(args...)
	} else {

		stmt := a.sqlStatement.ToStatementSQL(params)
		log.RootLogger().Debugf("Executing statement: %v", stmt)
		// stmt := a.sqlStatement.String()

		// rows = a.db.Exec(a.sqlStatement.ToStatementSQL(params), "abc", "test")
		// sqlStatement := `INSERT INTO readings_numeric(id, created, gatewayid, deviceid, resourceid, value) VALUES($1, $2, $3, $4, $5, $6)`
		// INSERT INTO users (age, email, first_name, last_name)
		// VALUES ($1, $2, $3, $4)`

		// 		sqlStatement := `
		// INSERT INTO users (age, email, first_name, last_name)
		// VALUES ($1, $2, $3, $4)`
		// sargs := a.sqlStatement.GetStatementArgs(params)
		// log.RootLogger().Infof("GetStatementArgs: %v", sargs...)

		// results, _ = a.db.Exec(stmt, "abc", "2021-05-26 19:24:21", "gate1", "device1", "res1", 100)

		results, _ = a.db.Exec(stmt)

	}
	if err != nil {
		return nil, err
	}

	// // defer rows.Close()

	// var results interface{}

	// if a.labeledResults {
	// 	results, err = getLabeledResults(a.dbHelper, rows)
	// } else {
	// 	results, err = getResults(a.dbHelper, rows)
	// }

	return results, nil
}

func getLabeledResults(dbHelper util.DbHelper, rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {

		values := make([]interface{}, len(columnTypes))
		for i := range values {
			values[i] = dbHelper.GetScanType(columnTypes[i])
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		resMap := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			resMap[column] = *(values[i].(*interface{}))
		}

		//todo do we need to do column mapping

		results = append(results, resMap)
	}

	return results, rows.Err()
}

func getResults(dbHelper util.DbHelper, rows *sql.Rows) ([][]interface{}, error) {

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var results [][]interface{}

	for rows.Next() {

		values := make([]interface{}, len(columnTypes))
		for i := range values {
			values[i] = dbHelper.GetScanType(columnTypes[i])
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		results = append(results, values)
	}

	return results, rows.Err()
}

//todo move to shared connection
func getConnection(s *Settings) (*sql.DB, error) {

	log.RootLogger().Debugf("getConnection DataSourceName:  %v", s.DataSourceName)

	db, err := sql.Open(s.DriverName, s.DataSourceName)
	if err != nil {
		return nil, err
	}

	if s.MaxOpenConns > 0 {
		db.SetMaxOpenConns(s.MaxOpenConns)
	}

	if s.MaxIdleConns != 2 {
		db.SetMaxIdleConns(s.MaxIdleConns)
	}

	return db, err
}
