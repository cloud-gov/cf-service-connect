package service

import (
	"fmt"

	"github.com/18F/cf-service-connect/models"
)

type pSQL struct{}

func (p pSQL) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("psql", "postgres")
}

// https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
func (p pSQL) GetConnectionUri(localPort int, creds models.Credentials) string {
	return fmt.Sprintf("postgresql://%s:%s@localhost:%d/%s", creds.GetUsername(), creds.GetPassword(), localPort, creds.GetDBName())
}

func (p pSQL) HasRepl() bool {
	return true
}

func (p pSQL) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		Cmd: "psql",
		Args: []string{
			// http://www.starkandwayne.com/blog/using-a-postgres-uri-with-psql/
			p.GetConnectionUri(localPort, creds),
		},
	}
}

// PSQL is the service singleton.
var PSQL = pSQL{}
