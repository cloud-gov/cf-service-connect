package connector

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/18F/cf-db-connect/models"

	"code.cloudfoundry.org/cli/plugin"
)

func Connect(cliConnection plugin.CliConnection, appName, serviceInstanceName string) (err error) {
	fmt.Println("Finding the service instance details...")

	service, err := cliConnection.GetService(serviceInstanceName)
	if err != nil {
		return
	}

	serviceName := service.ServiceOffering.Name
	planName := service.ServicePlan.Name

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

	serviceKeyCreds, err := getCreds(cliConnection, service.Guid, serviceKeyID)
	if err != nil {
		return
	}

	fmt.Println("Setting up SSH tunnel...")
	localPort, cmd, err := createSSHTunnel(serviceKeyCreds, appName)
	if err != nil {
		return
	}
	// TODO check if command failed

	// TODO ensure it works with Ctrl-C (exit early signal)

	if isMySQLService(serviceName, planName) {
		fmt.Println("Connecting to MySQL...")
		err = launchMySQL(localPort, serviceKeyCreds)
		if err != nil {
			return
		}
	} else if isPSQLService(serviceName, planName) {
		fmt.Println("Connecting to Postgres...")
		err = launchPSQL(localPort, serviceKeyCreds)
		if err != nil {
			return
		}
	} else {
		err = errors.New(fmt.Sprintf("Unsupported service. Service Name '%s' Plan Name '%s'. File an issue at https://github.com/18F/cf-db-connect/issues/new", serviceName, planName))
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

func createSSHTunnel(serviceKeyCreds models.Credentials, appName string) (localPort int, cmd *exec.Cmd, err error) {
	localPort = getAvailablePort()
	cmd = exec.Command("cf", "ssh", "-N", "-L", fmt.Sprintf("%d:%s:%s", localPort, serviceKeyCreds.Host, serviceKeyCreds.Port), appName)
	err = cmd.Start()
	return
}

func launchMySQL(localPort int, serviceKeyCreds models.Credentials) error {
	fmt.Printf("%+v\n", serviceKeyCreds)
	return startShell("mysql", []string{"-u", serviceKeyCreds.Username, "-h", "0", "-p" + serviceKeyCreds.Password, "-D", serviceKeyCreds.DBName, "-P", strconv.Itoa(localPort)})
}

func launchPSQL(localPort int, serviceKeyCreds models.Credentials) error {
	os.Setenv("PGPASSWORD", serviceKeyCreds.Password)
	return startShell("psql", []string{"-h", "localhost", "-p", fmt.Sprintf("%d", localPort), serviceKeyCreds.DBName, serviceKeyCreds.Username})
}

// derived from http://technosophos.com/2014/07/11/start-an-interactive-shell-from-within-go.html
func startShell(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait until user exits the shell
	return cmd.Run()
}

func getAvailablePort() int {
	// TODO find one that's available
	return 63306
}
func generateServiceKeyID() string {
	// TODO find one that's available, or randomize
	return "DB_CONNECT"
}

func isMySQLService(serviceName, planName string) bool {
	return isServiceType(serviceName, planName, "mysql")
}

func isPSQLService(serviceName, planName string) bool {
	return isServiceType(serviceName, planName, "psql", "postgres")
}

func isServiceType(serviceName, planName string, items ...string) bool {
	for _, item := range items {
		if strings.Contains(serviceName, item) || strings.Contains(planName, item) {
			return true
		}
	}
	return false
}
