package models

import (
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type ServiceInstance struct {
	GUID    string
	Name    string
	Service string
	Plan    string
}

func (si *ServiceInstance) IsMySQLService() bool {
	return si.isServiceType("mysql")
}

func (si *ServiceInstance) IsPSQLService() bool {
	return si.isServiceType("psql", "postgres")
}

func (si *ServiceInstance) isServiceType(items ...string) bool {
	for _, item := range items {
		if strings.Contains(si.Service, item) || strings.Contains(si.Plan, item) {
			return true
		}
	}
	return false
}

func FetchServiceInstance(cliConnection plugin.CliConnection, name string) (si ServiceInstance, err error) {
	srv, err := cliConnection.GetService(name)
	if err != nil {
		return
	}

	// serviceName := service.ServiceOffering.Name
	// planName := service.ServicePlan.Name
	si = ServiceInstance{
		GUID:    srv.Guid,
		Service: srv.ServiceOffering.Name,
		Plan:    srv.ServicePlan.Name,
	}
	return
}
