package models

import "encoding/json"

type serviceKeyResponse struct {
	Resources []serviceKeyResource `json:"resources"`
}

type serviceKeyResource struct {
	Entity struct {
		Credentials credentialsJSON `json:"credentials"`
	} `json:"entity"`
}

type Credentials interface {
	GetDBName() string
	GetHost() string
	GetUsername() string
	GetPassword() string
	GetPort() string
}

// http://stackoverflow.com/a/28035946/358804
type credentialsJSON struct {
	// these groups of fields should be interchangeable
	DBName string `json:"db_name"`
	Name   string `json:"name"`

	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	///////////////////////////////////////////////////

	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func (c credentialsJSON) GetDBName() string {
	if c.Name != "" {
		return c.Name
	}
	return c.DBName
}

func (c credentialsJSON) GetHost() string {
	if c.Host != "" {
		return c.Host
	}
	return c.Hostname
}

func (c credentialsJSON) GetUsername() string {
	return c.Username
}

func (c credentialsJSON) GetPassword() string {
	return c.Password
}

func (c credentialsJSON) GetPort() string {
	return c.Port
}

func CredentialsFromJSON(body string) (creds Credentials, err error) {
	serviceKeyResponse := serviceKeyResponse{}
	err = json.Unmarshal([]byte(body), &serviceKeyResponse)
	if err != nil {
		return
	}
	creds = serviceKeyResponse.Resources[0].Entity.Credentials

	return
}
