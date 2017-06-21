package config

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type SqlDB *sqlx.DB

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// DSN returns the Data Source Name
func (ci MySQLInfo) DSN() string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username + ":" +
		ci.Password + "@tcp(" +
		ci.Hostname + ":" +
		fmt.Sprintf("%d", ci.Port) + ")/" +
		ci.Name + ci.Parameter
}

// NewMySQL Connect to the database
func NewMySQL() (SqlDB, error) {
	var err error
	var sql *sqlx.DB
	// Connect to MySQL
	if sql, err = sqlx.Connect("mysql", MySQL.DSN()); err != nil {
		// logger.Fatalf("SQL Driver Error: %s", err.Error())
	}
	sql = sql.Unsafe()

	// Check if is alive
	return sql, sql.Ping()
}
