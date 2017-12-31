package service

import (
	"fmt"
	"strconv"

	"github.com/18F/cf-service-connect/models"
)

type mySQL struct{}

func (p mySQL) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("mysql")
}

// https://dev.mysql.com/doc/mysql-shell-excerpt/5.7/en/mysql-shell-connection-using-uri.html
func (p mySQL) GetConnectionUri(localPort int, creds models.Credentials) string {
	return fmt.Sprintf("mysql://%s:%s@localhost:%d/%s", creds.GetUsername(), creds.GetPassword(), localPort, creds.GetDBName())
}

func (p mySQL) HasRepl() bool {
	return true
}

func (p mySQL) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		Cmd: "mysql",
		Args: []string{
			"-u", creds.GetUsername(),
			"-h", "0",
			"-p" + creds.GetPassword(),
			"-D", creds.GetDBName(),
			"-P", strconv.Itoa(localPort),
		},
	}
}

// MySQL is the service singleton.
var MySQL = mySQL{}
