package launcher

import (
	"github.com/18F/cf-service-connect/logger"
	"os/exec"
	"strings"
)

func execute(name string, args ...string) *exec.Cmd {
	logger.Debugf("%s %s\n", name, strings.Join(args, " "))
	return exec.Command(name, args...)
}
