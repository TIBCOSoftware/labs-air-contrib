<!--
title: SQL Insert
weight: 4616
-->

# SQLInert Database Activity 
This activity provides your flogo application database insert statements. 


## Installation

```bash
flogo install github.com/TIBCOSoftware/labs-air-contrib/activity/sqlquery
```

## Configuration

### Settings:
| Name               | Type   | Description
|:---                | :---   | :---    
| dbType             | string | The type of database (mysql, oracle, postgres, sqlite, sqlserver) - **REQUIRED**         
| driverName         | string | The database driver name - **REQUIRED**
| dataSourceName     | string | The database DataSource name - **REQUIRED**
| maxOpenConnections | int    | Max open connections (default is unlimited)
| maxIdleConnections | int    | Max idle connections (default is 2)
| query              | string | The SQL select query - **REQUIRED**
| disablePrepared    | bool   | Disable prepared statement usage
| labeledResults     | bool   | Return results labeled by column name

### Input:
| Name   | Type | Description
|:---    | :--- | :---    
| params | map  |  The query parameters

### Output:
| Name        | Type  | Description
|:---         | :---  | :---    
| columnNames | array |  The names of the result columns
| results     | array |  The results

## Examples

### Insert
Simple query that gets all items with ID less than 10, retrieves all the columns.  In order to use *mysql*, you have to import the driver by adding `github.com/go-sql-driver/mysql` to 
the app imports section.  See [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) for more information on the driver.
```json
{
  "id": "dbquery",
  "name": "DbQuery",
  "activity": {
    "ref": "github.com/project-flogo/fc/activity/sqlquery",
    "settings": {
      "dbType": "mysql",
      "driverName": "mysql",
      "dataSourceName": "username:password@tcp(host:port)/dbName",
      "query": "select * from test where ID < 10"
    }
  }
}
```


Connectin to a database

db, err := sql.Open("postgres", dataSourceName)

connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
db, err := sql.Open("postgres", connStr)


Postgres connection

"host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"


Executing insert



result, err := db.ExecContext(ctx,
	"INSERT INTO users (name, age) VALUES ($1, $2)",
	"gopher",
	27,
)

Modifying app import

import (
	"database/sql"

	_ "github.com/lib/pq"
)

```
### Supported Drivers

- MySQL: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- Oracle: [github.com/mattn/go-oci8](https://github.com/mattn/go-oci8)
- Postgres: [github.com/lib/pq](https://github.com/lib/pq) 
- SQLite: [github.com/mattn/go-sqlite3]( https://github.com/mattn/go-sqlite3)
- SQLServer: [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb)
