package connector

import (
	"fmt"

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

	serviceKey := models.NewServiceKey(serviceInstance)

	// clean up existing service key, if present
	serviceKey.Delete(cliConnection)

	err = serviceKey.Create(cliConnection)
	if err != nil {
		return
	}
	defer serviceKey.Delete(cliConnection)

	creds, err := serviceKey.GetCreds(cliConnection)
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
