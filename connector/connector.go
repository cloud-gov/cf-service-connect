package connector

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/18F/cf-service-connect/launcher"
	"github.com/18F/cf-service-connect/models"
	"github.com/18F/cf-service-connect/service"

	"code.cloudfoundry.org/cli/plugin"
)

// Options are the structured representation of the command-line flags/arguments.
type Options struct {
	AppName             string
	ServiceInstanceName string
	ConnectClient       bool
}

const manualConnectInstructions = `Skipping call to client CLI. Connection information:

Host: localhost
Port: {{.Port}}
Username: {{.User}}
Password: {{.Pass}}
Name: {{.Name}}

{{if .HasRepl }}To connect:

    {{.LaunchCmd}}

{{end}}Leave this terminal open while you want to use the SSH tunnel. Press Control-C to stop.
`

type localConnectionData struct {
	Port      int
	User      string
	Pass      string
	Name      string
	HasRepl   bool
	LaunchCmd service.LaunchCmd
}

func manualConnect(srv service.Service, tunnel launcher.SSHTunnel, creds models.Credentials) (err error) {
	launchCmd := srv.GetLaunchCmd(tunnel.LocalPort, creds)
	connectionData := localConnectionData{
		Port:      tunnel.LocalPort,
		User:      creds.GetUsername(),
		Pass:      creds.GetPassword(),
		Name:      creds.GetDBName(),
		HasRepl:   srv.HasRepl(),
		LaunchCmd: launchCmd,
	}

	tmpl, err := template.New("").Parse(manualConnectInstructions)
	if err != nil {
		return
	}
	err = tmpl.Execute(os.Stdout, connectionData)
	if err != nil {
		return
	}

	// wait for a Control-C
	tunnel.Wait()

	return
}

func handleClient(
	options Options,
	tunnel launcher.SSHTunnel,
	si models.ServiceInstance,
	creds models.Credentials,
) error {
	srv := service.GetService(si)
	if srv == service.UnknownService {
		fmt.Printf("Unable to find matching client for service '%s' with plan '%s'.\n", si.Service, si.Plan)
	}

	if options.ConnectClient {
		if srv.HasRepl() {
			cmd := srv.GetLaunchCmd(tunnel.LocalPort, creds)
			// check whether the executable is available
			_, err := exec.LookPath(cmd.Cmd)
			if err == nil {
				fmt.Println("Connecting client...")
				return cmd.Exec()
			} else {
				fmt.Printf("Executable `%s` not found.\n", cmd.Cmd)
			}
		}

		fmt.Println("Falling back to `-no-client` behavior.")
	}

	return manualConnect(srv, tunnel, creds)
}

// Connect performs the primary action of the plugin: providing an SSH tunnel and launching the appropriate client, if desired.
func Connect(cliConnection plugin.CliConnection, options Options) (err error) {
	fmt.Println("Finding the service instance details...")

	serviceInstance, err := models.FetchServiceInstance(cliConnection, options.ServiceInstanceName)
	if err != nil {
		return
	}

	serviceKey := models.NewServiceKey(serviceInstance)

	// clean up existing service key, if present
	serviceKey.Delete(cliConnection)

	err = serviceKey.Create(cliConnection)
	if err != nil {
		return
	}
	defer serviceKey.Delete(cliConnection)

	creds, err := serviceKey.GetCreds(cliConnection)
	if err != nil {
		return
	}

	fmt.Println("Setting up SSH tunnel...")
	tunnel := launcher.NewSSHTunnel(creds, options.AppName)
	err = tunnel.Open()
	if err != nil {
		return
	}
	defer tunnel.Close()

	err = handleClient(options, tunnel, serviceInstance, creds)
	return
}
