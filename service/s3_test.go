package service

import (
	"testing"

	"github.com/18F/cf-service-connect/models"
	"github.com/stretchr/testify/assert"
)

type s3MatchTest struct {
	serviceName string
	planName    string
	expected    bool
}

func TestS3Match(t *testing.T) {
	tests := []s3MatchTest{
		{
			"s3",
			"shared",
			true,
		},
		{
			"aws",
			"s3",
			true,
		},
		{
			"mysql",
			"shared",
			false,
		},
	}

	for _, test := range tests {
		serviceInstance := models.ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		result := S3.Match(serviceInstance)
		assert.Equal(t, result, test.expected)
	}
}
