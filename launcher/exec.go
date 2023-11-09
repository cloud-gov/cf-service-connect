package launcher

import (
	"os/exec"
	"strings"

	"github.com/18F/cf-service-connect/logger"
)

func execute(name string, args ...string) *exec.Cmd {
	logger.Debugf("%s %s\n", name, strings.Join(args, " "))
	return exec.Command(name, args...) // #nosec
}
