package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/cloud-gov/cf-service-connect/connector"

	"code.cloudfoundry.org/cli/plugin"
	cf_client "github.com/cloudfoundry/go-cfclient/v3/client"
	cf_client_config "github.com/cloudfoundry/go-cfclient/v3/config"
)

const subcommand = "connect-to-service"

// ServiceConnectPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type ServiceConnectPlugin struct{}

func (c *ServiceConnectPlugin) parseOptions(args []string) (options connector.Options, err error) {
	metadata := c.GetMetadata()
	command := metadata.Commands[0]
	flags := flag.NewFlagSet(command.Name, flag.ExitOnError)
	option := "no-client"
	noClient := flags.Bool(option, false, command.UsageDetails.Options[option])

	err = flags.Parse(args[1:])
	if err != nil {
		return
	}

	nonFlagArgs := flags.Args()
	if len(nonFlagArgs) != 2 {
		err = errors.New("Wrong number of arguments")
		return
	}

	options = connector.Options{
		AppName:             nonFlagArgs[0],
		ServiceInstanceName: nonFlagArgs[1],
		ConnectClient:       !(*noClient),
	}
	return
}

// Run is the entry point when the core CLI is invoking a command defined
// by the plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second paramter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
func (c *ServiceConnectPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// check to ensure it's the right subcommand, not others like CLI-MESSAGE-UNINSTALL
	if args[0] != subcommand {
		return
	}

	hasApiEndpoint, err := cliConnection.HasAPIEndpoint()
	if err != nil || !hasApiEndpoint {
		err = fmt.Errorf("no API endpoint set")
		return
	}

	loggedIn, err := cliConnection.IsLoggedIn()
	if err != nil {
		return
	}
	if !loggedIn {
		err = fmt.Errorf("please log in to search for apps")
		return
	}

	cfc, err := createCfClient()
	if err != nil {
		return
	}

	opts, err := c.parseOptions(args)
	if err != nil {
		log.Fatalln(err)
	}

	// connect has plugin.CliConnection already instantiated in c.Run.
	err = connector.Connect(cliConnection, opts)
	if err != nil {
		log.Fatalln(err)
	}
}

func createCfClient() (*cf_client.Client, error) {
	cfg, err := cf_client_config.NewFromCFHome()
	if err != nil {
		return &cf_client.Client{}, err
	}

	cfc, err := cf_client.New(cfg)
	if err != nil {
		return &cf_client.Client{}, err
	}

	return cfc, nil
}

// GetMetadata returns the plugin information for the CLI to consume.
func (c *ServiceConnectPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "ServiceConnect",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 9,
			Build: 9,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 15,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     subcommand,
				HelpText: "Open a shell that's connected to a database service instance",
				UsageDetails: plugin.Usage{
					Usage: "\n   cf " + subcommand + " [-no-client] <app_name> <service_instance_name>",
					Options: map[string]string{
						"no-client": "If this param is passed, the CLI client for the service won't be started, and the connection information will be printed to the console. Useful for connecting to the service through a GUI.",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(ServiceConnectPlugin))
}
