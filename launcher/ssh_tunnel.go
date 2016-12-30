package launcher

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/18F/cf-db-connect/models"
)

func getAvailablePort() int {
	// TODO find one that's available
	return 63306
}

type SSHTunnel struct {
	LocalPort int
	cmd       *exec.Cmd
}

func (t *SSHTunnel) Open() (err error) {
	// Start the (long-running) SSH tunnel command, and ensure that it doesn't fail. Because of how Process management in Go works, this needs to happen asynchronously.
	// https://groups.google.com/d/msg/golang-nuts/XviHC3bJF8s/PUOYzcsmvwMJ

	errChan := make(chan error)

	go func() {
		errChan <- t.cmd.Run()
	}()

	// if the tunnel will fail to be created, it *should* be done in this time
	time.Sleep(4 * time.Second)

	select {
	default:
		// success!
	case e := <-errChan:
		// SSH tunnel failed
		if e == nil {
			err = errors.New("SSH tunnel command exited early, without error")
		} else {
			err = e
		}
	}

	return
}

func (t *SSHTunnel) Close() error {
	return t.cmd.Process.Kill()
}

func NewSSHTunnel(creds models.Credentials, appName string) SSHTunnel {
	localPort := getAvailablePort()

	cmd := exec.Command("cf", "ssh", "-N", "-L", fmt.Sprintf("%d:%s:%s", localPort, creds.Host, creds.Port), appName)
	// should only print in the case of an issue
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return SSHTunnel{
		LocalPort: localPort,
		cmd:       cmd,
	}
}
