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
		assert.Nil(t, err)

		assert.Equal(t, creds.GetHost(), test.expectedHost)
		assert.Equal(t, creds.GetDBName(), test.expectedDBName)
	}
}
