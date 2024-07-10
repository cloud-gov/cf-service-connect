package service

import (
	"strconv"

	"github.com/cloud-gov/cf-service-connect/launcher"
	"github.com/cloud-gov/cf-service-connect/models"
)

type redis struct{}

func (p redis) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("redis")
}

func (p redis) Launch(localPort int, creds models.Credentials) error {
	return launcher.StartShell("redis-cli", getRedisLaunchFlags(localPort, creds))
}

func getRedisLaunchFlags(localPort int, creds models.Credentials) []string {
	return []string{
		"--tls",
		"-p", strconv.Itoa(localPort),
		"-a", creds.GetPassword(),
	}
}

// Redis is the service singleton.
var Redis = redis{}
