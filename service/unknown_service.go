package service

import (
	"github.com/18F/cf-service-connect/models"
)

type unknownService struct{}

func (p unknownService) Match(si models.ServiceInstance) bool {
	return true
}

func (p unknownService) GetConnectionUri(localPort int, creds models.Credentials) string {
	return "unknown"
}

func (p unknownService) HasRepl() bool {
	return false
}

func (p unknownService) GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd {
	return LaunchCmd{
		// no-op
		// https://stackoverflow.com/a/12405621/358804
		Cmd: ":",
	}
}

// UnknownService is the service singleton.
var UnknownService = unknownService{}
