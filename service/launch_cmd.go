package service

import (
	"os"
	"strings"

	"github.com/18F/cf-service-connect/launcher"
)

type LaunchCmd struct {
	Envs map[string]string
	Cmd  string
	Args []string
}

func (lc LaunchCmd) Exec() error {
	for envVar, val := range lc.Envs {
		err := os.Setenv(envVar, val)
		if err != nil {
			return err
		}
	}

	return launcher.StartShell(lc.Cmd, lc.Args)
}

func (lc LaunchCmd) String() string {
	result := ""
	for envVar, val := range lc.Envs {
		result += envVar + "=" + val + " "
	}
	result += lc.Cmd
	if len(lc.Args) > 0 {
		result += " " + strings.Join(lc.Args, " ")
	}
	return result
}
