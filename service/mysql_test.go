package service

import (
	"testing"

	"github.com/18F/cf-db-connect/models"
	"github.com/stretchr/testify/assert"
)

type mySQLMatchTest struct {
	serviceName string
	planName    string
	expected    bool
}

func TestMySQLMatch(t *testing.T) {
	tests := []mySQLMatchTest{
		{
			"mysql",
			"shared",
			true,
		},
		{
			"somedb",
			"shared-mysql",
			true,
		},
		{
			"aws",
			"rds",
			false,
		},
		{
			"psql",
			"shared",
			false,
		},
	}

	for _, test := range tests {
		serviceInstance := models.ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		result := MySQL.Match(serviceInstance)
		assert.Equal(t, result, test.expected)
	}
}
