package models

import (
	"encoding/json"
	"strconv"
)

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
	GetPort() (int, error)
	GetRegion() string
}

// http://stackoverflow.com/a/28035946/358804
type credentialsJSON struct {
	// these groups of fields should be interchangeable
	DBName string `json:"db_name"`
	Dbname string `json:"dbname"`
	Bucket string `json:"bucket"`
	Name   string `json:"name"`

	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	HostName string `json:"host_name"`

	Username    string `json:"username"`
	UserName    string `json:"user_name"`
	AccessKeyId string `json:"access_key_id"`
	User        string `json:"user"`

	Password        string `json:"password"`
	SecretAccessKey string `json:"secret_access_key"`
	Pass            string `json:"pass"`

	Region string `json:"region"`
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
	if c.Bucket != "" {
		return c.Bucket
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
	if c.AccessKeyId != "" {
		return c.AccessKeyId
	}
	return c.User
}

func (c credentialsJSON) GetPassword() string {
	if c.Pass != "" {
		return c.Pass
	}
	if c.SecretAccessKey != "" {
		return c.SecretAccessKey
	}
	return c.Password
}

func (c credentialsJSON) GetPort() (int, error) {
	portStr := c.Port.String()
	if portStr == "" {
		// TODO think of a better way to handle this
		return 0, nil
	}
	return strconv.Atoi(portStr)
}

func (c credentialsJSON) GetRegion() string {
	return c.Region
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
