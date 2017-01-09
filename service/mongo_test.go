package service

import (
	"testing"

	"github.com/18F/cf-service-connect/models"
	"github.com/stretchr/testify/assert"
)

type mongoDBMatchTest struct {
	serviceName string
	planName    string
	expected    bool
}

func TestMongoDBMatch(t *testing.T) {
	tests := []mongoDBMatchTest{
		{
			"mongo",
			"shared",
			true,
		},
		{
			"mongodb",
			"shared",
			true,
		},
		{
			"somedb",
			"shared-mongo",
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
		result := MongoDB.Match(serviceInstance)
		assert.Equal(t, result, test.expected)
	}
}
