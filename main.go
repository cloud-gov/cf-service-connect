package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"encoding/json"
	"net/url"
	"os/exec"

	"code.cloudfoundry.org/cli/plugin"
)

const SUBCOMMAND = "connect-to-db"

// DBConnectPlugin is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type DBConnectPlugin struct{}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second paramter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.
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

	fmt.Println("Finding the service instance details...")

	service, err := cliConnection.GetService(serviceInstanceName)
	if err != nil {
		log.Fatalln(err)
	}

	serviceName := service.ServiceOffering.Name
	planName := service.ServicePlan.Name

	serviceKeyID := generateServiceKeyID()
	_, err = cliConnection.CliCommandWithoutTerminalOutput("create-service-key", serviceInstanceName, serviceKeyID)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_, err := cliConnection.CliCommandWithoutTerminalOutput("delete-service-key", "-f", serviceInstanceName, serviceKeyID)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	serviceKeyCreds, err := getCreds(cliConnection, service.Guid, serviceKeyID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Setting up SSH tunnel...")
	localPort, cmd, err := createSSHTunnel(serviceKeyCreds, appName)
	if err != nil {
		log.Fatalln(err)
	}
	// TODO check if command failed

	// TODO ensure it works with Ctrl-C (exit early signal)

	if isMySQLService(serviceName, planName) {
		fmt.Println("Connecting to MySQL...")
		err = launchMySQL(localPort, serviceKeyCreds)
		if err != nil {
			log.Fatalln(err)
		}
	} else if isPSQLService(serviceName, planName) {
		fmt.Println("Connecting to Postgres...")
		err = launchPSQL(localPort, serviceKeyCreds)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalf("Unsupported service. Service Name '%s' Plan Name '%s'. File an issue at https://github.com/18F/cf-db-connect/issues/new", serviceName, planName)
	}

	// TODO defer
	if err := cmd.Process.Kill(); err != nil {
		log.Println(err)
	}
}

func getCreds(cliConnection plugin.CliConnection, serviceGUID, serviceKeyID string) (creds Credentials, err error) {
	serviceKeyAPI := fmt.Sprintf("/v2/service_instances/%s/service_keys?q=name%%3A%s", serviceGUID, url.QueryEscape(serviceKeyID))
	bodyLines, err := cliConnection.CliCommandWithoutTerminalOutput("curl", serviceKeyAPI)
	if err != nil {
		return
	}
	body := strings.Join(bodyLines, "")
	serviceKeyResponse := ServiceKeyResponse{}
	err = json.Unmarshal([]byte(body), &serviceKeyResponse)
	if err != nil {
		return
	}

	creds = serviceKeyResponse.Resources[0].Entity.Credentials
	return
}

func createSSHTunnel(serviceKeyCreds Credentials, appName string) (localPort int, cmd *exec.Cmd, err error) {
	localPort = getAvailablePort()
	cmd = exec.Command("cf", "ssh", "-N", "-L", fmt.Sprintf("%d:%s:%s", localPort, serviceKeyCreds.Host, serviceKeyCreds.Port), appName)
	err = cmd.Start()
	return
}

func launchMySQL(localPort int, serviceKeyCreds Credentials) error {
	fmt.Printf("%+v\n", serviceKeyCreds)
	return startShell("mysql", []string{"-u", serviceKeyCreds.Username, "-h", "0", "-p" + serviceKeyCreds.Password, "-D", serviceKeyCreds.DBName, "-P", strconv.Itoa(localPort)})
}

func launchPSQL(localPort int, serviceKeyCreds Credentials) error {
	os.Setenv("PGPASSWORD", serviceKeyCreds.Password)
	return startShell("psql", []string{"-h", "localhost", "-p", fmt.Sprintf("%d", localPort), serviceKeyCreds.DBName, serviceKeyCreds.Username})
}

// derived from http://technosophos.com/2014/07/11/start-an-interactive-shell-from-within-go.html
func startShell(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait until user exits the shell
	return cmd.Run()
}

type ServiceKeyResponse struct {
	Resources []ServiceKeyResource `json:"resources"`
}

type ServiceKeyResource struct {
	Entity struct {
		Credentials Credentials `json:"credentials"`
	} `json:"entity"`
}

type Credentials struct {
	DBName   string `json:"db_name"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func getAvailablePort() int {
	// TODO find one that's available
	return 63306
}
func generateServiceKeyID() string {
	// TODO find one that's available, or randomize
	return "DB_CONNECT"
}

func isMySQLService(serviceName, planName string) bool {
	return isServiceType(serviceName, planName, "mysql")
}

func isPSQLService(serviceName, planName string) bool {
	return isServiceType(serviceName, planName, "psql", "postgres")
}

func isServiceType(serviceName, planName string, items ...string) bool {
	for _, item := range items {
		if strings.Contains(serviceName, item) || strings.Contains(planName, item) {
			return true
		}
	}
	return false
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as seperate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
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
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     SUBCOMMAND,
				HelpText: "Basic plugin command's help text",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: SUBCOMMAND + "\n   cf " + SUBCOMMAND + " <app_name> <service_instance_name>",
				},
			},
		},
	}
}

// Unlike most Go programs, the `Main()` function will not be used to run all of the
// commands provided in your plugin. Main will be used to initialize the plugin
// process, as well as any dependencies you might require for your
// plugin.
func main() {
	// Any initialization for your plugin can be handled here
	//
	// Note: to run the plugin.Start method, we pass in a pointer to the struct
	// implementing the interface defined at "code.cloudfoundry.org/cli/plugin/plugin.go"
	//
	// Note: The plugin's main() method is invoked at install time to collect
	// metadata. The plugin will exit 0 and the Run([]string) method will not be
	// invoked.
	plugin.Start(new(DBConnectPlugin))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}
