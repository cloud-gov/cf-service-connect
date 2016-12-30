package models

import (
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

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

func (sk *ServiceKey) GetCreds(cliConnection plugin.CliConnection) (creds Credentials, err error) {
	serviceKeyAPI := fmt.Sprintf("/v2/service_instances/%s/service_keys?q=name%%3A%s", sk.Instance.GUID, url.QueryEscape(sk.ID))
	bodyLines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", serviceKeyAPI)
	if err != nil {
		return
	}

	body := strings.Join(bodyLines, "")
	creds, err = CredentialsFromJSON(body)
	return
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
