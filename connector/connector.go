package connector

import (
	"fmt"
	"os"
	"text/template"

	"github.com/18F/cf-db-connect/launcher"
	"github.com/18F/cf-db-connect/models"
	"github.com/18F/cf-db-connect/service"

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

Leave this terminal open while you want to use the SSH tunnel. Press Control-C to stop.
`

type localConnectionData struct {
	Port int
	User string
	Pass string
	Name string
}

func manualConnect(tunnel launcher.SSHTunnel, creds models.Credentials) (err error) {
	connectionData := localConnectionData{
		Port: tunnel.LocalPort,
		User: creds.GetUsername(),
		Pass: creds.GetPassword(),
		Name: creds.GetDBName(),
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
	if options.ConnectClient {
		srv, found := service.GetService(si)
		if found {
			fmt.Println("Connecting client...")
			return srv.Launch(tunnel.LocalPort, creds)
		}

		fmt.Printf("Unable to find matching client for service '%s' with plan '%s'. Falling back to `-no-client` behavior.\n", si.Service, si.Plan)
	}

	return manualConnect(tunnel, creds)
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
