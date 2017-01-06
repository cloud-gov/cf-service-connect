package launcher

import "os"

// derived from http://technosophos.com/2014/07/11/start-an-interactive-shell-from-within-go.html
func StartShell(name string, args []string) error {
	cmd := execute(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait until user exits the shell
	return cmd.Run()
}
