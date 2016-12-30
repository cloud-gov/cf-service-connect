package launcher

import (
	"fmt"
	"os/exec"

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

func (t *SSHTunnel) Open() error {
	return t.cmd.Start()
}

func (t *SSHTunnel) Close() error {
	return t.cmd.Process.Kill()
}

func NewSSHTunnel(serviceKeyCreds models.Credentials, appName string) SSHTunnel {
	localPort := getAvailablePort()
	cmd := exec.Command("cf", "ssh", "-N", "-L", fmt.Sprintf("%d:%s:%s", localPort, serviceKeyCreds.Host, serviceKeyCreds.Port), appName)

	return SSHTunnel{
		LocalPort: localPort,
		cmd:       cmd,
	}
}
