package connector

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/cloud-gov/cf-service-connect/api"
	"github.com/cloud-gov/cf-service-connect/launcher"
	"github.com/cloud-gov/cf-service-connect/models"
	"github.com/cloud-gov/cf-service-connect/service"
)

// Options are the structured representation of the command-line
// flags/arguments.
type Options struct {
	AppName             string
	ServiceInstanceName string
	ConnectClient       bool
}

// startedAppState is the CF app state that indicates the operator wants the app
// running. It is a desired state, not a guarantee that an instance exists.
const startedAppState = "STARTED"

// The plugin tunnels through the app's first web instance, matching `cf ssh`'s
// default. These are named constants rather than literals so the intent is
// visible at the call site.
const (
	sshProcessType   = "web"
	sshInstanceIndex = 0
)

const manualConnectInstructions = `Skipping call to client CLI. Connection information:

Host: localhost
Port: {{.Port}}
Username: {{.User}}
Password: {{.Pass}}
Name: {{.Name}}
`

const disconnectInstructions = `
Leave this terminal open while you want to use the SSH tunnel. Press Control-C to stop.`

type localConnectionData struct {
	Port int
	User string
	Pass string
	Name string
}

func getConnectionInstructions(creds models.Credentials) string {
	instructions := manualConnectInstructions
	if creds.IsPostgresDatabase() {
		pgConnectionEnvVars := `
You can set these environment variables to connect to your PostgreSQL database:

export PGHOST=localhost
export PGPORT={{.Port}}
export PGUSER="{{.User}}"
export PGPASSWORD="{{.Pass}}"
export PGDATABASE="{{.Name}}"
`
		instructions = instructions + pgConnectionEnvVars
	}
	instructions = instructions + disconnectInstructions
	return instructions
}

func manualConnect(tunnel *launcher.SSHTunnel, creds models.Credentials) error {
	connectionData := localConnectionData{
		Port: tunnel.LocalPort,
		User: creds.GetUsername(),
		Pass: creds.GetPassword(),
		Name: creds.GetDBName(),
	}

	instructions := getConnectionInstructions(creds)
	tmpl, err := template.New("").Parse(instructions)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(os.Stdout, connectionData); err != nil {
		return err
	}

	// Wait for either a Control-C or the tunnel failing on its own. Previously
	// only the tunnel was watched, so a tunnel that died left the user staring
	// at instructions for a connection that no longer worked.
	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupted)

	tunnelClosed := make(chan error, 1)
	go func() {
		tunnelClosed <- tunnel.Wait()
	}()

	select {
	case <-interrupted:
		fmt.Println("\nClosing the SSH tunnel.")
		return nil
	case err := <-tunnelClosed:
		return err
	}
}

func handleClient(
	options Options,
	tunnel *launcher.SSHTunnel,
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

// Connect performs the primary action of the plugin: providing an SSH tunnel
// and launching the appropriate client, if desired.
func Connect(cliConnection plugin.CliConnection, options Options) error {
	return connect(context.Background(), api.NewConnection(cliConnection), options)
}

func connect(ctx context.Context, conn api.Connection, options Options) (err error) {
	launcher.WarnIfCFBinaryNameSet()

	client, err := api.NewClient(conn)
	if err != nil {
		return err
	}

	fmt.Println("Finding the service instance details...")
	instance, err := client.GetServiceInstance(ctx, options.ServiceInstanceName)
	if err != nil {
		return err
	}
	serviceInstance := models.ServiceInstance{
		GUID:    instance.GUID,
		Name:    instance.Name,
		Service: instance.Offering,
		Plan:    instance.Plan,
	}

	app, err := client.GetApp(ctx, options.AppName)
	if err != nil {
		return err
	}
	if app.State != startedAppState {
		return fmt.Errorf("app %q is not started (state: %s); an SSH tunnel needs a running app instance.\nStart it with: cf start %s",
			app.Name, app.State, app.Name)
	}

	// Resolve everything the SSH session needs *before* creating a service key.
	// All of these can fail for reasons the user must fix (SSH disabled, no
	// running instance, no SSH endpoint), and there is no point provisioning
	// credentials we are then going to throw away.
	target, err := sshTarget(ctx, client, app)
	if err != nil {
		return err
	}

	serviceKey := models.NewServiceKey(serviceInstance)

	// Clean up a key left behind by an earlier interrupted run. Unlike before,
	// a failure here is reported: it usually means the key exists but cannot be
	// deleted, in which case the create below would fail with a less obvious
	// "already exists" error.
	if err := deleteExistingServiceKey(ctx, client, serviceKey); err != nil {
		return err
	}

	fmt.Println("Creating the service key...")
	serviceKey.GUID, err = client.CreateServiceKey(ctx, serviceKey.Instance.GUID, serviceKey.Name)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("Deleting the service key...")
		if deleteErr := client.DeleteServiceKey(ctx, serviceKey.GUID); deleteErr != nil {
			// Report but do not mask the primary error: a leaked key is worth
			// telling the user about, since they may need to remove it by hand
			// before the next run.
			fmt.Fprintf(os.Stderr,
				"Warning: could not delete the temporary service key %q: %v\nRemove it with: cf delete-service-key %s %s\n",
				serviceKey.Name, deleteErr, serviceKey.Instance.Name, serviceKey.Name)
			if err == nil {
				err = deleteErr
			}
		}
	}()

	rawCreds, err := client.GetServiceKeyCredentials(ctx, serviceKey.GUID)
	if err != nil {
		return err
	}
	creds, err := models.CredentialsFromMap(rawCreds)
	if err != nil {
		return err
	}

	fmt.Println("Setting up SSH tunnel...")
	tunnel := launcher.NewSSHTunnel(creds, target)
	if err := tunnel.Open(); err != nil {
		return err
	}
	defer func() {
		if closeErr := tunnel.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return handleClient(options, tunnel, serviceInstance, creds)
}

func deleteExistingServiceKey(ctx context.Context, client *api.Client, serviceKey models.ServiceKey) error {
	guid, found, err := client.FindServiceKey(ctx, serviceKey.Instance.GUID, serviceKey.Name)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	fmt.Printf("Removing the service key %q left over from a previous run...\n", serviceKey.Name)
	if err := client.DeleteServiceKey(ctx, guid); err != nil {
		return fmt.Errorf("could not remove the existing service key %q: %w", serviceKey.Name, err)
	}
	return nil
}

// sshTarget assembles everything needed to open an SSH session to the app: the
// proxy endpoint from the CF API root document, the process instance to target,
// and a one-time passcode from UAA.
func sshTarget(ctx context.Context, client *api.Client, app api.App) (launcher.SSHTarget, error) {
	endpoint, err := client.GetSSHEndpoint(ctx)
	if err != nil {
		return launcher.SSHTarget{}, err
	}

	if err := client.CheckSSHEnabled(ctx, app); err != nil {
		return launcher.SSHTarget{}, err
	}

	// The app being STARTED does not mean an instance is running, and the SSH
	// username needs the process GUID rather than the app GUID.
	process, err := client.GetSSHProcess(ctx, app, sshProcessType, sshInstanceIndex)
	if err != nil {
		return launcher.SSHTarget{}, err
	}

	passcode, err := client.SSHPasscode(ctx)
	if err != nil {
		return launcher.SSHTarget{}, err
	}
	if passcode == "" {
		return launcher.SSHTarget{}, errors.New("UAA returned an empty SSH passcode")
	}

	return launcher.SSHTarget{
		Address:            endpoint.Address,
		WebSocketURL:       endpoint.WebSocketURL,
		HostKeyFingerprint: endpoint.HostKeyFingerprint,
		User:               process.SSHUsername(),
		Passcode:           passcode,
		TLSConfig:          client.TLSConfig(),
	}, nil
}
