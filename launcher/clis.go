package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/18F/cf-db-connect/models"
)

// derived from http://technosophos.com/2014/07/11/start-an-interactive-shell-from-within-go.html
func startShell(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait until user exits the shell
	return cmd.Run()
}

func LaunchMySQL(localPort int, creds models.Credentials) error {
	fmt.Printf("%+v\n", creds)
	return startShell("mysql", []string{
		"-u", creds.Username,
		"-h", "0",
		"-p" + creds.Password,
		"-D", creds.DBName,
		"-P", strconv.Itoa(localPort),
	})
}

func LaunchPSQL(localPort int, creds models.Credentials) error {
	os.Setenv("PGPASSWORD", creds.Password)
	return startShell("psql", []string{
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", localPort),
		creds.DBName,
		creds.Username,
	})
}
