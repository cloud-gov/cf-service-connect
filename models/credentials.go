package models

import "encoding/json"

type ServiceKeyResponse struct {
	Resources []ServiceKeyResource `json:"resources"`
}

type ServiceKeyResource struct {
	Entity struct {
		Credentials Credentials `json:"credentials"`
	} `json:"entity"`
}

type Credentials struct {
	DBName   string `json:"db_name"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func CredentialsFromJSON(body string) (creds Credentials, err error) {
	serviceKeyResponse := ServiceKeyResponse{}
	err = json.Unmarshal([]byte(body), &serviceKeyResponse)
	if err != nil {
		return
	}
	creds = serviceKeyResponse.Resources[0].Entity.Credentials
	return
}
