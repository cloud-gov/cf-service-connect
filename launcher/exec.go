package launcher

import (
	"fmt"
	"github.com/18F/cf-db-connect/utils"
	"os"
	"os/exec"
	"strings"
)

func execute(name string, args ...string) *exec.Cmd {
	utils.Debugf("%s %s\n", name, strings.Join(args, " "))
	return exec.Command(name, args...)
}
