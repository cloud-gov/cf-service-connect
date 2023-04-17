package service

import (
	"fmt"

	"github.com/18F/cf-service-connect/models"
)

type pSQL struct{}

func (p pSQL) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("psql", "postgres")
}

func (p pSQL) HasRepl() bool {
	return true
}

func (p pSQL) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		Envs: map[string]string{
			"PGPASSWORD": creds.GetPassword(),
		},
		Cmd: "psql",
		Args: []string{
			"-h", "localhost",
			"-p", fmt.Sprintf("%d", localPort),
			creds.GetDBName(),
			creds.GetUsername(),
		},
	}
}

// PSQL is the service singleton.
var PSQL = pSQL{}
