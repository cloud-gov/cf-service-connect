package main

import (
	"fmt"
	"log"
	"os"

	"github.com/18F/cf-db-connect/connector"

	"code.cloudfoundry.org/cli/plugin"
)

const SUBCOMMAND = "connect-to-db"

// DBConnectPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type DBConnectPlugin struct{}

// Run is the entry point when the core CLI is invoking a command defined
// by the plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second paramter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
func (c *DBConnectPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// check to ensure it's the right subcommand, not others like CLI-MESSAGE-UNINSTALL
	if args[0] != SUBCOMMAND {
		return
	}

	if len(args) != 3 {
		metadata := c.GetMetadata()
		fmt.Println("Wrong number of arguments. Usage:")
		fmt.Println(metadata.Commands[0].UsageDetails.Usage)
		os.Exit(1)
	}

	appName := args[1]
	serviceInstanceName := args[2]
	err := connector.Connect(cliConnection, appName, serviceInstanceName)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *DBConnectPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DBConnect",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 15,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     SUBCOMMAND,
				HelpText: "Basic plugin command's help text",
				UsageDetails: plugin.Usage{
					Usage: SUBCOMMAND + "\n   cf " + SUBCOMMAND + " <app_name> <service_instance_name>",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DBConnectPlugin))
}
