package service

import (
	"fmt"

	"github.com/18F/cf-service-connect/models"
)

type elasticsearch struct{}

func (p elasticsearch) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("elasticsearch")
}

func (p elasticsearch) GetConnectionUri(localPort int, creds models.Credentials) string {
	return fmt.Sprintf("http://%s:%s@localhost:%d", creds.GetUsername(), creds.GetPassword(), localPort)
}

func (p elasticsearch) HasRepl() bool {
	return false
}

func (p elasticsearch) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		Cmd: "curl",
		Args: []string{
			p.GetConnectionUri(localPort, creds),
		},
	}
}

// Elasticsearch is the service singleton.
var Elasticsearch = elasticsearch{}
