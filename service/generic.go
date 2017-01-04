package service

import "github.com/18F/cf-db-connect/models"

type Service interface {
	Match(si models.ServiceInstance) bool
}

var SERVICES = []Service{
	MySQL{},
	PSQL{},
}

func GetService(si models.ServiceInstance) (Service, bool) {
	for _, potentialService := range SERVICES {
		if potentialService.Match(si) {
			return potentialService, true
		}
	}

	return nil, false
}
