package service

import (
	"strconv"

	"github.com/18F/cf-service-connect/models"
)

type mongoDB struct{}

func (p mongoDB) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("mongo")
}

func (p mongoDB) HasRepl() bool {
	return true
}

func (p mongoDB) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		Cmd: "mongo",
		Args: []string{
			"-u", creds.GetUsername(),
			"-p", creds.GetPassword(),
			"--port", strconv.Itoa(localPort),
			creds.GetDBName(),
		},
	}
}

// MongoDB is the service singleton.
var MongoDB = mongoDB{}
