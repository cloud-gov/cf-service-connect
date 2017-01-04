package service

import (
	"testing"

	"github.com/18F/cf-db-connect/models"
	"github.com/stretchr/testify/assert"
)

type getServiceTest struct {
	serviceName   string
	planName      string
	expectFound   bool
	expectService Service
}

func TestGetService(t *testing.T) {
	tests := []getServiceTest{
		{
			"psql",
			"shared",
			true,
			PSQL{},
		},
		{
			"mysql",
			"shared",
			true,
			MySQL{},
		},
		{
			"other",
			"service",
			false,
			nil,
		},
	}

	for _, test := range tests {
		serviceInstance := models.ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		srv, found := GetService(serviceInstance)
		assert.Equal(t, found, test.expectFound)
		if test.expectFound {
			assert.Equal(t, srv, test.expectService)
		}
	}

}
