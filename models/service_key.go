package models

import "code.cloudfoundry.org/cli/plugin"

type ServiceKey struct {
	Instance ServiceInstance
	ID       string
}

func (sk *ServiceKey) Create(conn plugin.CliConnection) error {
	_, err := conn.CliCommandWithoutTerminalOutput("create-service-key", sk.Instance.Name, sk.ID)
	return err
}

func (sk *ServiceKey) Delete(conn plugin.CliConnection) error {
	_, err := conn.CliCommandWithoutTerminalOutput("delete-service-key", "-f", sk.Instance.Name, sk.ID)
	return err
}

func generateServiceKeyID() string {
	// TODO find one that's available, or randomize
	return "DB_CONNECT"
}

func NewServiceKey(instance ServiceInstance) ServiceKey {
	return ServiceKey{
		Instance: instance,
		ID:       generateServiceKeyID(),
	}
}
