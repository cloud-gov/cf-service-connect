package connector

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/18F/cf-db-connect/launcher"
	"github.com/18F/cf-db-connect/models"

	"code.cloudfoundry.org/cli/plugin"
)

func Connect(cliConnection plugin.CliConnection, appName, serviceInstanceName string) (err error) {
	fmt.Println("Finding the service instance details...")

	serviceInstance, err := models.FetchServiceInstance(cliConnection, serviceInstanceName)
	if err != nil {
		return
	}

	serviceKey := models.NewServiceKey(serviceInstanceName)

	// clean up existing service key, if present
	serviceKey.Delete(cliConnection)

	err = serviceKey.Create(cliConnection)
	if err != nil {
		return
	}
	defer serviceKey.Delete(cliConnection)

	creds, err := getCreds(cliConnection, serviceInstance.GUID, serviceKey.ID)
	if err != nil {
		return
	}

	fmt.Println("Setting up SSH tunnel...")
	tunnel := launcher.NewSSHTunnel(creds, appName)
	err = tunnel.Open()
	if err != nil {
		return
	}
	defer tunnel.Close()

	err = launcher.LaunchDBCLI(serviceInstance, tunnel, creds)
	return
}

func getCreds(cliConnection plugin.CliConnection, serviceGUID, serviceKeyID string) (creds models.Credentials, err error) {
	serviceKeyAPI := fmt.Sprintf("/v2/service_instances/%s/service_keys?q=name%%3A%s", serviceGUID, url.QueryEscape(serviceKeyID))
	bodyLines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", serviceKeyAPI)
	if err != nil {
		return
	}

	body := strings.Join(bodyLines, "")
	creds, err = models.CredentialsFromJSON(body)
	return
}
