package connector

import (
	"errors"
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

	serviceKeyID := generateServiceKeyID()

	// clean up existing service key, if present
	deleteServiceKey(cliConnection, serviceInstanceName, serviceKeyID)

	_, err = cliConnection.CliCommandWithoutTerminalOutput("create-service-key", serviceInstanceName, serviceKeyID)
	if err != nil {
		return
	}
	defer func() {
		err := deleteServiceKey(cliConnection, serviceInstanceName, serviceKeyID)
		if err != nil {
			return
		}
	}()

	serviceKeyCreds, err := getCreds(cliConnection, serviceInstance.GUID, serviceKeyID)
	if err != nil {
		return
	}

	fmt.Println("Setting up SSH tunnel...")
	localPort, cmd, err := launcher.CreateSSHTunnel(serviceKeyCreds, appName)
	if err != nil {
		return
	}
	// TODO check if command failed

	// TODO ensure it works with Ctrl-C (exit early signal)

	if isMySQLService(serviceInstance) {
		fmt.Println("Connecting to MySQL...")
		err = launcher.LaunchMySQL(localPort, serviceKeyCreds)
		if err != nil {
			return
		}
	} else if isPSQLService(serviceInstance) {
		fmt.Println("Connecting to Postgres...")
		err = launcher.LaunchPSQL(localPort, serviceKeyCreds)
		if err != nil {
			return
		}
	} else {
		msg := fmt.Sprintf("Unsupported service. Service Name '%s' Plan Name '%s'. File an issue at https://github.com/18F/cf-db-connect/issues/new", serviceInstance.Service, serviceInstance.Plan)
		err = errors.New(msg)
		return
	}

	// TODO defer
	err = cmd.Process.Kill()
	return
}

func deleteServiceKey(conn plugin.CliConnection, serviceInstanceName, serviceKeyID string) error {
	_, err := conn.CliCommandWithoutTerminalOutput("delete-service-key", "-f", serviceInstanceName, serviceKeyID)
	return err
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

func generateServiceKeyID() string {
	// TODO find one that's available, or randomize
	return "DB_CONNECT"
}

func isMySQLService(si models.ServiceInstance) bool {
	return isServiceType(si.Service, si.Plan, "mysql")
}

func isPSQLService(si models.ServiceInstance) bool {
	return isServiceType(si.Service, si.Plan, "psql", "postgres")
}

func isServiceType(serviceName, planName string, items ...string) bool {
	for _, item := range items {
		if strings.Contains(serviceName, item) || strings.Contains(planName, item) {
			return true
		}
	}
	return false
}
