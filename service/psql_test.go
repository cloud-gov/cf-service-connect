package service

import (
	"testing"

	"github.com/18F/cf-db-connect/models"
	"github.com/stretchr/testify/assert"
)

type pSQLMatchTest struct {
	serviceName string
	planName    string
	expected    bool
}

func TestPSQLMatch(t *testing.T) {
	tests := []pSQLMatchTest{
		{
			"psql",
			"shared",
			true,
		},
		{
			"postgres",
			"shared",
			true,
		},
		{
			"somedb",
			"psql",
			true,
		},
		{
			"somedb",
			"postgres",
			true,
		},
		{
			"aws",
			"rds",
			false,
		},
		{
			"mysql",
			"shared",
			false,
		},
	}

	pSQL := PSQL{}
	for _, test := range tests {
		serviceInstance := models.ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		result := pSQL.Match(serviceInstance)
		assert.Equal(t, result, test.expected)
	}
}
