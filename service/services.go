package service

import "github.com/18F/cf-service-connect/models"

type Service interface {
	Match(si models.ServiceInstance) bool
	GetConnectionUri(localPort int, creds models.Credentials) string
	HasRepl() bool
	GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd
}

var services = []Service{
	Elasticsearch,
	MongoDB,
	MySQL,
	PSQL,
	Redis,
}

func GetService(si models.ServiceInstance) Service {
	for _, potentialService := range services {
		if potentialService.Match(si) {
			return potentialService
		}
	}

	return UnknownService
}
