package service

import (
	"testing"

	"github.com/18F/cf-service-connect/models"
	"github.com/stretchr/testify/assert"
)

type redisMatchTest struct {
	serviceName string
	planName    string
	expected    bool
}

func TestRedisMatch(t *testing.T) {
	tests := []redisMatchTest{
		{
			"redis",
			"shared",
			true,
		},
		{
			"somedb",
			"shared-redis",
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
		result := Redis.Match(serviceInstance)
		assert.Equal(t, result, test.expected)
	}
}
