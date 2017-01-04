package connector

import (
	"fmt"
	"os"
	"text/template"

	"github.com/18F/cf-db-connect/launcher"
	"github.com/18F/cf-db-connect/models"

	"code.cloudfoundry.org/cli/plugin"
)

type Options struct {
	AppName             string
	ServiceInstanceName string
	ConnectClient       bool
}

func manualConnect(tunnel launcher.SSHTunnel, creds models.Credentials) (err error) {
	rawTemplate := `Skipping call to client CLI. Connection information:

Host: localhost
Port: {{.Port}}
Username: {{.User}}
Password: {{.Pass}}
Name: {{.Name}}

Leave this terminal open while you want to use the SSH tunnel. Press Control-C to stop.
`

	// http://julianyap.com/2013/09/23/using-anonymous-structs-to-pass-data-to-templates-in-golang.html
	connectionData := struct {
		Port int
		User string
		Pass string
		Name string
	}{
		tunnel.LocalPort,
		creds.GetUsername(),
		creds.GetPassword(),
		creds.GetDBName(),
	}

	tmpl, err := template.New("").Parse(rawTemplate)
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

	if options.ConnectClient {
		err = launcher.LaunchDBCLI(serviceInstance, tunnel, creds)
		return
	} else {
		err = manualConnect(tunnel, creds)
		return
	}

	return
}
