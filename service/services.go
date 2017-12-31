package service

import "github.com/18F/cf-service-connect/models"

type Service interface {
	Match(si models.ServiceInstance) bool
	GetLaunchCmd(localPort int, creds models.Credentials) LaunchCmd
}

var services = []Service{
	MongoDB,
	MySQL,
	PSQL,
	Redis,
}

func GetService(si models.ServiceInstance) (Service, bool) {
	for _, potentialService := range services {
		if potentialService.Match(si) {
			return potentialService, true
		}
	}

	return nil, false
}
