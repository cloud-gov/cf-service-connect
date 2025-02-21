package api

import (
	"fmt"
	"net/url"

	plugin_models "code.cloudfoundry.org/cli/plugin/models"
)

type CFClient struct {
	connection    Connection
	CCAPIEndpoint string
	AuthToken     string
	AppId         string
	AppName       string
	IsSSLDisabled bool
}

type Connection interface {
	ApiEndpoint() (string, error)
	HasSpace() (bool, error)
	IsLoggedIn() (bool, error)
	AccessToken() (string, error)
	GetCurrentSpace() (plugin_models.Space, error)
	IsSSLDisabled() (bool, error)
}

func NewCFClient(connection Connection) (*CFClient, error) {

	ccAPIEndpoint, err := connection.ApiEndpoint()
	if err != nil {
		return nil, err
	}
	if ccAPIEndpoint == "" {
		return nil, fmt.Errorf("No Cloud Foundry api endpoint set")
	}

	isSSLDisabled, err := connection.IsSSLDisabled()
	if err != nil {
		return nil, err
	}

	client := &CFClient{
		connection:    connection,
		CCAPIEndpoint: ccAPIEndpoint,
		IsSSLDisabled: isSSLDisabled,
	}

	return client, nil

}

func (client *CFClient) Configure(appName string) error {

	if connected, err := client.connection.IsLoggedIn(); !connected {
		if err != nil {
			return err
		}
		return fmt.Errorf("You must be logged in %s first.", client.CCAPIEndpoint)
	}

	if hasSpace, err := client.connection.HasSpace(); !hasSpace {
		if err != nil {
			return err
		}
		return fmt.Errorf("No org and space targeted, use 'cf target -o ORG -s SPACE' to target an org and space")
	}

	currentSpace, err := client.connection.GetCurrentSpace()
	if err != nil {
		return err
	}

	authToken, err := client.connection.AccessToken()
	if err != nil {
		return err
	}

	ccAPIURL, err := url.Parse(client.CCAPIEndpoint)
	if err != nil {
		return err
	}

	cfAPIClient, err := NewCFAPIClient(ccAPIURL, authToken, client.IsSSLDisabled)
	if err != nil {
		return err
	}

	appGUID, err := cfAPIClient.GetAppGUID(appName, currentSpace.Guid)
	if err != nil {
		return err
	}

	client.AuthToken = authToken
	client.AppId = appGUID
	client.AppName = appName
	return nil

}
