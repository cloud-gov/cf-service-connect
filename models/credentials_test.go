package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// credentialsMap parses a JSON object the way api.Client hands it over: as the
// decoded "credentials" object from
// GET /v3/service_credential_bindings/:guid/details.
func credentialsMap(t *testing.T, body string) map[string]any {
	t.Helper()

	var raw map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &raw))
	return raw
}

type credentialsTest struct {
	name           string
	credsJSON      string
	expectedHost   string
	expectedPort   string
	expectedDBName string
	expectedUser   string
	expectedPass   string
	expectedUri    string
}

// The broker field names below are all in active use; the plugin has to accept
// every variant.
func credentialsTests() []credentialsTest {
	return []credentialsTest{
		{
			name: "host/db_name/username/password",
			credsJSON: `{
				"host": "host.com",
				"port": "5432",
				"db_name": "name",
				"username": "user",
				"password": "pass",
				"uri": "uri"
			}`,
			expectedHost:   "host.com",
			expectedPort:   "5432",
			expectedDBName: "name",
			expectedUser:   "user",
			expectedPass:   "pass",
			expectedUri:    "uri",
		},
		{
			name: "hostname/name/user/pass",
			credsJSON: `{
				"hostname": "host.com",
				"port": "5432",
				"name": "name",
				"user": "user",
				"pass": "pass"
			}`,
			expectedHost:   "host.com",
			expectedPort:   "5432",
			expectedDBName: "name",
			expectedUser:   "user",
			expectedPass:   "pass",
		},
		{
			name: "host_name/user_name",
			credsJSON: `{
				"host_name": "host.com",
				"port": "5432",
				"name": "name",
				"user_name": "user",
				"password": "pass"
			}`,
			expectedHost:   "host.com",
			expectedPort:   "5432",
			expectedDBName: "name",
			expectedUser:   "user",
			expectedPass:   "pass",
		},
		{
			// CAPI serves the port as a JSON number for the aws-rds broker, so
			// it must not be assumed to be a string.
			name: "numeric port and dbname",
			credsJSON: `{
				"host_name": "host.com",
				"port": 5432,
				"dbname": "name",
				"user_name": "user",
				"password": "pass"
			}`,
			expectedHost:   "host.com",
			expectedPort:   "5432",
			expectedDBName: "name",
			expectedUser:   "user",
			expectedPass:   "pass",
		},
	}
}

func TestCredentialsFromMap(t *testing.T) {
	for _, test := range credentialsTests() {
		t.Run(test.name, func(t *testing.T) {
			creds, err := CredentialsFromMap(credentialsMap(t, test.credsJSON))
			require.NoError(t, err)

			assert.Equal(t, test.expectedHost, creds.GetHost())
			assert.Equal(t, test.expectedPort, creds.GetPort())
			assert.Equal(t, test.expectedDBName, creds.GetDBName())
			assert.Equal(t, test.expectedUser, creds.GetUsername())
			assert.Equal(t, test.expectedPass, creds.GetPassword())
		})
	}
}

// A v2-shaped payload has no usable host at the top level, because v2 nested the
// credentials under resources[0].entity.credentials. It must be rejected with a
// clear message rather than silently yielding empty connection details.
func TestCredentialsFromMapRejectsAV2Envelope(t *testing.T) {
	v2Body := `{
		"resources": [
			{
				"entity": {
					"credentials": {
						"host": "host.com",
						"port": "5432",
						"db_name": "name",
						"username": "user",
						"password": "pass"
					}
				}
			}
		]
	}`

	_, err := CredentialsFromMap(credentialsMap(t, v2Body))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no host")
}

// Regression test: the previous implementation indexed resources[0] without a
// bounds check, so any JSON response lacking that key panicked. A CAPI v3 error
// body is exactly such a response.
func TestCredentialsFromMapDoesNotPanicOnAnErrorBody(t *testing.T) {
	errorBody := `{"errors":[{"detail":"Unknown request","title":"CF-NotFound","code":10000}]}`

	assert.NotPanics(t, func() {
		_, err := CredentialsFromMap(credentialsMap(t, errorBody))
		assert.Error(t, err)
	})
}

func TestCredentialsFromMapRejectsEmpty(t *testing.T) {
	_, err := CredentialsFromMap(map[string]any{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no credentials")
}

// A service returning credentials without a port cannot be tunnelled to, so this
// must fail early rather than at connection time.
func TestCredentialsRejectMissingPort(t *testing.T) {
	_, err := CredentialsFromMap(credentialsMap(t, `{"host": "host.com", "username": "u", "password": "p"}`))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no port")
}

func TestCredentialsRejectMissingHost(t *testing.T) {
	_, err := CredentialsFromMap(credentialsMap(t, `{"port": 5432, "username": "u", "password": "p"}`))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no host")
}

// An unexpected type for a credential field must be an error, not a panic or a
// silently wrong value.
func TestCredentialsFromMapRejectsAWrongFieldType(t *testing.T) {
	_, err := CredentialsFromMap(map[string]any{
		"host": []string{"host.com"},
		"port": 5432,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not parse")
}
