package launcher

import (
	"os/exec"
	"strings"

	"github.com/cloud-gov/cf-service-connect/logger"
)

func execute(name string, args ...string) *exec.Cmd {
	logger.Debugf("%s %s\n", name, strings.Join(args, " "))
	return exec.Command(name, args...)
}
