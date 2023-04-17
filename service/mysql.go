package service

import (
	"strconv"

	"github.com/18F/cf-service-connect/models"
)

type mySQL struct{}

func (p mySQL) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("mysql")
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
