package config

import (
	"fmt"
)

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
