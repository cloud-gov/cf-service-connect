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

func CredentialsFromJSON(body string) (creds Credentials, err error) {
	response := serviceKeyResponse{}
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return
	}
	creds = response.Resources[0].Entity.Credentials

	return
}
