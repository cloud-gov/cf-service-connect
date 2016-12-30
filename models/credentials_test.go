package models

import (
	"fmt"
	"testing"
)

const JSONWrap = `{
	"resources": [
		{
			"entity": {
				"credentials": %s
			}
		}
	]
}`

type credentialsFromJSONTest struct {
	entityJSON     string
	expectedDBName string
	expectedHost   string
}

func TestCredentialsFromJSON(t *testing.T) {
	tests := []credentialsFromJSONTest{
		{
			`{
				"db_name": "name",
				"host": "host.com",
				"password": "pass",
				"port": "5432",
				"username": "user"
			}`,
			"name",
			"host.com",
		},
		{
			`{
				"name": "name",
				"hostname": "host.com",
				"password": "pass",
				"port": "5432",
				"username": "user"
			}`,
			"name",
			"host.com",
		},
	}

	for _, test := range tests {
		fullJSON := fmt.Sprintf(JSONWrap, test.entityJSON)
		creds, err := CredentialsFromJSON(fullJSON)
		if err != nil {
			t.Error(err)
		}
		if creds.GetHost() != test.expectedHost {
			t.Errorf("Expected: %v. Actual: %v.", test.expectedHost, creds.GetHost())
		}
		if creds.GetDBName() != test.expectedDBName {
			t.Errorf("Expected: %v. Actual: %v.", test.expectedDBName, creds.GetDBName())
		}
	}
}
