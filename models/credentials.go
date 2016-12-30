package models

import "encoding/json"

type ServiceKeyResponse struct {
	Resources []ServiceKeyResource `json:"resources"`
}

type ServiceKeyResource struct {
	Entity struct {
		Credentials CredentialsJSON `json:"credentials"`
	} `json:"entity"`
}

// http://stackoverflow.com/a/28035946/358804
type CredentialsJSON struct {
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

func (c *CredentialsJSON) GetDBName() string {
	if c.Name != "" {
		return c.Name
	}
	return c.DBName
}

func (c *CredentialsJSON) GetHost() string {
	if c.Host != "" {
		return c.Host
	}
	return c.Hostname
}

type Credentials struct {
	DBName   string
	Host     string
	Username string
	Password string
	Port     string
}

func CredentialsFromJSON(body string) (creds Credentials, err error) {
	serviceKeyResponse := ServiceKeyResponse{}
	err = json.Unmarshal([]byte(body), &serviceKeyResponse)
	if err != nil {
		return
	}
	jsonCreds := serviceKeyResponse.Resources[0].Entity.Credentials

	creds = Credentials{
		DBName:   jsonCreds.GetDBName(),
		Host:     jsonCreds.GetHost(),
		Username: jsonCreds.Username,
		Password: jsonCreds.Password,
		Port:     jsonCreds.Port,
	}

	return
}
