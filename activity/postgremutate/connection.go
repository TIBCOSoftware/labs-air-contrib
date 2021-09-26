package connection

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "Postgres.connection")
var factory = &PgFactory{}

// Settings for postgres
type Settings struct {
	DatabaseURL string `md:"databaseURL"`
	Host        string `md:"host,required"`
	Port        int    `md:"port,required"`
	User        string `md:"user,required"`
	Password    string `md:"password,required"`
	DbName      string `md:"databaseName,required"`
	SSLMode     string `md:"sslmode"`
	// Value interface{} `md:"value"`
		
}



func init() {
	if os.Getenv(log.EnvKeyLogLevel) == "DEBUG" {
		 logCache.DebugEnabled()
	}

	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

// PgFactory for postgres connection
type PgFactory struct {
}

// Type PgFactory
func (*PgFactory) Type() string {
	return "Postgres"
}

// NewManager PgFactory
func (*PgFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {

	
	sharedConn := &PgSharedConfigManager{
		
	}
	var err error
	

	s := &Settings{}

	err = metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}

	cHost := s.Host
	if cHost == "" {
		return nil,errors.New("Required Parameter Host Name is missing")
	}
	
	cPort := s.Port
	if cPort == 0 {
		return nil,errors.New("Required Parameter Port is missing")
	}
	cDbName := s.DbName
	if cDbName == "" {
		return nil,errors.New("Required Parameter Database name is missing")
	}
	cUser := s.User
	if cUser == "" {
		return nil,errors.New("Required Parameter User is missing")
	}
	cPassword := s.Password
	if cPassword == "" {
		return nil,errors.New("Required Parameter Password is missing")
	}
	


	conninfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cHost, cPort, cUser, cPassword, cDbName)

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return nil,fmt.Errorf("Could not open connection to database %s, %s", cDbName, err.Error())
	}
	sharedConn.db = db

	err = db.Ping()
	if err != nil {
		return nil,err
	}

	return sharedConn, nil
}

// PgSharedConfigManager details
type PgSharedConfigManager struct {
	name      string
	db          *sql.DB
}
// Type PgSharedConfigManager details
func (p *PgSharedConfigManager) Type() string {
	return "Postgres"
}
// GetConnection PgSharedConfigManager details 
func (p *PgSharedConfigManager) GetConnection() interface{} {
	return p.db
}

// ReleaseConnection PgSharedConfigManager details
func (p *PgSharedConfigManager) ReleaseConnection(connection interface{}) {

}

// Start PgSharedConfigManager details
func (p *PgSharedConfigManager) Start() error {
	return nil
}
// Stop PgSharedConfigManager details
func (p *PgSharedConfigManager) Stop() error {
	logCache.Debug("Cleaning up DB")
     p.db.Close()

	return nil
}



