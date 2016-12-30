package models

import "code.cloudfoundry.org/cli/plugin"

type ServiceInstance struct {
	GUID    string
	Service string
	Plan    string
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
