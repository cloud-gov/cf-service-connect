package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	expectedHost   string
	expectedPort   string
	expectedDBName string
	expectedUser   string
	expectedPass   string
}

func TestCredentialsFromJSON(t *testing.T) {
	tests := []credentialsFromJSONTest{
		{
			`{
				"host": "host.com",
				"port": "5432",
				"db_name": "name",
				"username": "user",
				"password": "pass"
			}`,
			"host.com",
			"5432",
			"name",
			"user",
			"pass",
		},
		{
			`{
				"hostname": "host.com",
				"port": "5432",
				"name": "name",
				"user": "user",
				"pass": "pass"
			}`,
			"host.com",
			"5432",
			"name",
			"user",
			"pass",
		},
		{
			`{
				"host_name": "host.com",
				"port": "5432",
				"name": "name",
				"user_name": "user",
				"password": "pass"
			}`,
			"host.com",
			"5432",
			"name",
			"user",
			"pass",
		},
		{
			`{
				"host_name": "host.com",
				"port": 5432,
				"dbname": "name",
				"user_name": "user",
				"password": "pass"
			}`,
			"host.com",
			"5432",
			"name",
			"user",
			"pass",
		},
	}

	for _, test := range tests {
		fullJSON := fmt.Sprintf(JSONWrap, test.entityJSON)
		creds, err := CredentialsFromJSON(fullJSON)
		assert.Nil(t, err)

		assert.Equal(t, creds.GetHost(), test.expectedHost)
		assert.Equal(t, creds.GetPort(), test.expectedPort)
		assert.Equal(t, creds.GetDBName(), test.expectedDBName)
		assert.Equal(t, creds.GetUsername(), test.expectedUser)
		assert.Equal(t, creds.GetPassword(), test.expectedPass)
	}
}
