package launcher

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/cloud-gov/cf-service-connect/models"
	"github.com/phayes/freeport"
)

func getAvailablePort() int {
	return freeport.GetPort()
}

// SSHTunnel represents an underlying SSH tunnel process. The struct should be created via NewSSHTunnel().
type SSHTunnel struct {
	LocalPort int
	cmd       *exec.Cmd
	errChan   chan error
}

// Open starts the underlying SSH tunnel process.
func (t *SSHTunnel) Open() (err error) {
	// Start the (long-running) SSH tunnel command, and ensure that it doesn't fail. Because of how Process management in Go works, this needs to happen asynchronously.
	// https://groups.google.com/d/msg/golang-nuts/XviHC3bJF8s/PUOYzcsmvwMJ

	go func() {
		t.errChan <- t.cmd.Run()
	}()

	// if the tunnel will fail to be created, it *should* be done in this time
	time.Sleep(6 * time.Second)

	select {
	default:
		fmt.Println("SSH tunnel created.") // we hope
	case e := <-t.errChan:
		// SSH tunnel failed
		if e == nil {
			err = errors.New("SSH tunnel command exited early, without error")
		} else {
			err = e
		}
	}

	return
}

// Wait will block until the SSH tunnel is closed, either intentionally or not.
func (t *SSHTunnel) Wait() error {
	return <-t.errChan
}

// Close will terminate the underlying SSH tunnel process.
func (t *SSHTunnel) Close() error {
	return t.cmd.Process.Kill()
}

// NewSSHTunnel prepares the underlying Cloud Foundry CLI process for creating an SSH tunnel to a service instance.
func NewSSHTunnel(creds models.Credentials, appName string) SSHTunnel {
	localPort := getAvailablePort()

	cmd := execute(
		getCFBinaryName(),
		"ssh",
		"-N",
		"-L", fmt.Sprintf("%d:%s:%s", localPort, creds.GetHost(), creds.GetPort()),
		appName,
	)
	// should only print in the case of an issue
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return SSHTunnel{
		LocalPort: localPort,
		cmd:       cmd,
		errChan:   make(chan error),
	}
}

func getCFBinaryName() string {
	cfBinaryName, exists := os.LookupEnv("CF_BINARY_NAME")
	if !exists {
		cfBinaryName = "cf"
	}
	return cfBinaryName
}
