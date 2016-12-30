package connector

import (
	"errors"
	"fmt"

	"github.com/18F/cf-db-connect/launcher"
	"github.com/18F/cf-db-connect/models"

	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
)

func Connect(cliConnection plugin.CliConnection, serviceInstanceName string) (err error) {
	fmt.Println("Finding the service instance details...")

	serviceInstance, err := models.FetchServiceInstance(cliConnection, serviceInstanceName)
	if err != nil {
		return
	}

	serviceKey := models.NewServiceKey(serviceInstance)

	fmt.Println("Cleaning up existing service keys...")
	serviceKey.Delete(cliConnection)

	fmt.Println("Creating new service key...")
	err = serviceKey.Create(cliConnection)
	if err != nil {
		return
	}
	defer serviceKey.Delete(cliConnection)

	fmt.Println("Retrieving service key credentials...")
	creds, err := serviceKey.GetCreds(cliConnection)
	if err != nil {
		return
	}

	fmt.Println("Determining app to bind to...")
	app, err := getApp(cliConnection)
	if err != nil {
		return
	}

	fmt.Println("Setting up SSH tunnel...")
	tunnel := launcher.NewSSHTunnel(creds, app.Name)
	err = tunnel.Open()
	if err != nil {
		return
	}
	defer tunnel.Close()

	err = launcher.LaunchDBCLI(serviceInstance, tunnel, creds)
	return
}

func getApp(conn plugin.CliConnection) (app plugin_models.GetAppsModel, err error) {
	apps, err := conn.GetApps()
	if err != nil {
		return
	}
	if len(apps) == 0 {
		err = errors.New("No apps in this space")
		return
	}

	app = apps[0]
	return
}
