package api

import (
	"context"
	"fmt"
	"net/url"

	cf_client "github.com/cloudfoundry/go-cfclient/v3/client"
	cf_client_config "github.com/cloudfoundry/go-cfclient/v3/config"
)

type CFAPIClient struct {
	client *cf_client.Client
}

func NewCFAPIClient(ccAPIEndpoint *url.URL, authToken string, skipTLSValidation bool) (*CFAPIClient, error) {
	// A refresh token is not provided by the CF CLI Plugin API and is not required as
	// "AccessToken() now provides a refreshed o-auth token.",
	// see https://github.com/cloudfoundry/cli/blob/main/plugin/plugin_examples/CHANGELOG.md#changes-in-v614
	refreshToken := ""

	cfClientConfigOptions := []cf_client_config.Option{
		cf_client_config.Token(authToken, refreshToken),
	}

	if skipTLSValidation {
		cfClientConfigOptions = append(cfClientConfigOptions, cf_client_config.SkipTLSValidation())
	}

	cfg, err := cf_client_config.New(ccAPIEndpoint.String(), cfClientConfigOptions...)
	if err != nil {
		return nil, err
	}
	cf, err := cf_client.New(cfg)
	if err != nil {
		return nil, err
	}
	return &CFAPIClient{client: cf}, nil
}

func (client *CFAPIClient) GetAppGUID(appName string, currentSpaceGUID string) (string, error) {
	appFilter := &cf_client.AppListOptions{
		Names:      cf_client.Filter{Values: []string{appName}},
		SpaceGUIDs: cf_client.Filter{Values: []string{currentSpaceGUID}},
	}
	apps, err := client.client.Applications.ListAll(context.Background(), appFilter)
	if err != nil {
		return "", err
	}

	if len(apps) == 0 {
		return "", fmt.Errorf("App '%s' not found", appName)
	}

	app := apps[0]
	return app.GUID, nil
}
