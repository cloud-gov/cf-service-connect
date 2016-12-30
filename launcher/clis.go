package launcher

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/18F/cf-db-connect/logger"
	"github.com/18F/cf-db-connect/models"
)

// derived from http://technosophos.com/2014/07/11/start-an-interactive-shell-from-within-go.html
func startShell(name string, args []string) error {
	cmd := execute(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait until user exits the shell
	return cmd.Run()
}

func LaunchMySQL(localPort int, creds models.Credentials) error {
	return startShell("mysql", []string{
		"-u", creds.GetUsername(),
		"-h", "0",
		"-p" + creds.GetPassword(),
		"-D", creds.GetDBName(),
		"-P", strconv.Itoa(localPort),
	})
}

func LaunchPSQL(localPort int, creds models.Credentials) error {
	os.Setenv("PGPASSWORD", creds.GetPassword())
	logger.Debugf("PGPASSWORD=%s ", creds.GetPassword())

	return startShell("psql", []string{
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", localPort),
		creds.GetDBName(),
		creds.GetUsername(),
	})
}

func LaunchDBCLI(serviceInstance models.ServiceInstance, tunnel SSHTunnel, creds models.Credentials) error {
	if serviceInstance.IsMySQLService() {
		fmt.Println("Connecting to MySQL...")
		return LaunchMySQL(tunnel.LocalPort, creds)
	} else if serviceInstance.IsPSQLService() {
		fmt.Println("Connecting to Postgres...")
		return LaunchPSQL(tunnel.LocalPort, creds)
	} else {
		msg := fmt.Sprintf("Unsupported service. Service Name '%s' Plan Name '%s'. File an issue at https://github.com/18F/cf-db-connect/issues/new", serviceInstance.Service, serviceInstance.Plan)
		return errors.New(msg)
	}
}
