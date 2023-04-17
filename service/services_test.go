package service

import (
	"testing"

	"github.com/18F/cf-service-connect/models"
	"github.com/stretchr/testify/assert"
)

type getServiceTest struct {
	serviceName   string
	planName      string
	expectService Service
}

func TestGetService(t *testing.T) {
	tests := []getServiceTest{
		{
			"psql",
			"shared",
			PSQL,
		},
		{
			"mysql",
			"shared",
			MySQL,
		},
		{
			"other",
			"service",
			UnknownService,
		},
	}

	for _, test := range tests {
		serviceInstance := models.ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		srv := GetService(serviceInstance)
		assert.Equal(t, srv, test.expectService)
	}

}
