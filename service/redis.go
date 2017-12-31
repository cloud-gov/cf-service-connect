package service

import (
	"fmt"
	"strconv"

	"github.com/18F/cf-service-connect/models"
)

type redis struct{}

func (p redis) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("redis")
}

// https://www.iana.org/assignments/uri-schemes/prov/redis
func (p redis) GetConnectionUri(localPort int, creds models.Credentials) string {
	return fmt.Sprintf("redis://%s:%s@localhost:%d/%s", creds.GetUsername(), creds.GetPassword(), localPort, creds.GetDBName())
}

func (p redis) HasRepl() bool {
	return true
}

func (p redis) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		Cmd: "redis-cli",
		Args: []string{
			"-p", strconv.Itoa(localPort),
			"-a", creds.GetPassword(),
		},
	}
}

// Redis is the service singleton.
var Redis = redis{}
