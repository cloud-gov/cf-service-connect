package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Credentials exposes the connection details of a service instance.
type Credentials interface {
	GetDBName() string
	GetHost() string
	GetUsername() string
	GetPassword() string
	GetPort() string
	IsPostgresDatabase() bool
}

// credentialsJSON accommodates the differing field names that service brokers
// use for the same concepts.
//
// http://stackoverflow.com/a/28035946/358804
type credentialsJSON struct {
	// these groups of fields should be interchangeable
	DBName string `json:"db_name"`
	Dbname string `json:"dbname"`
	Name   string `json:"name"`

	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	HostName string `json:"host_name"`

	Username string `json:"username"`
	UserName string `json:"user_name"`
	User     string `json:"user"`

	Password string `json:"password"`
	Pass     string `json:"pass"`
	///////////////////////////////////////////////////

	Uri string `json:"uri"`

	// can be an integer or a string
	// http://igorsobreira.com/2015/04/11/decoding-json-numbers-into-strings-in-go.html
	Port json.Number `json:"port"`
}

func (c credentialsJSON) GetDBName() string {
	if c.Name != "" {
		return c.Name
	}
	if c.Dbname != "" {
		return c.Dbname
	}
	return c.DBName
}

func (c credentialsJSON) GetHost() string {
	if c.Host != "" {
		return c.Host
	}
	if c.HostName != "" {
		return c.HostName
	}
	return c.Hostname
}

func (c credentialsJSON) GetUsername() string {
	if c.Username != "" {
		return c.Username
	}
	if c.UserName != "" {
		return c.UserName
	}
	return c.User
}

func (c credentialsJSON) GetPassword() string {
	if c.Pass != "" {
		return c.Pass
	}
	return c.Password
}

func (c credentialsJSON) GetPort() string {
	return c.Port.String()
}

func (c credentialsJSON) IsPostgresDatabase() bool {
	return strings.HasPrefix(c.Uri, "postgres")
}

// CredentialsFromMap converts the credentials object returned by
// GET /v3/service_credential_bindings/:guid/details into Credentials.
//
// CAPI v3 returns the credentials as a flat object, whereas v2 nested them under
// resources[0].entity.credentials.
func CredentialsFromMap(raw map[string]any) (Credentials, error) {
	if len(raw) == 0 {
		return nil, errors.New("the service key returned no credentials")
	}

	// Round-trip through JSON so that the interchangeable field-name handling
	// and json.Number port parsing above apply, rather than duplicating that
	// logic for map[string]any.
	encoded, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("could not re-encode the service key credentials: %w", err)
	}

	var creds credentialsJSON
	if err := json.Unmarshal(encoded, &creds); err != nil {
		return nil, fmt.Errorf("could not parse the service key credentials: %w", err)
	}

	// A broker returning credentials in an unrecognised shape would otherwise
	// surface much later as a confusing connection failure, so check the fields
	// the tunnel actually needs.
	if creds.GetHost() == "" {
		return nil, errors.New("the service key credentials contain no host; this service may not be supported")
	}
	if creds.GetPort() == "" {
		return nil, errors.New("the service key credentials contain no port; this service may not be supported")
	}

	return creds, nil
}
