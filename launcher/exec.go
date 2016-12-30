package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execWithDebug(name string, args ...string) *exec.Cmd {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("%s %s\n", name, strings.Join(args, " "))
	}

	return exec.Command(name, args...)
}
