package models

import (
	"testing"

	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/stretchr/testify/assert"
)

type isServiceTest struct {
	serviceName string
	planName    string
	result      bool
}

func TestIsMySQLService(t *testing.T) {
	tests := []isServiceTest{
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
		serviceInstance := ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		result := serviceInstance.IsMySQLService()
		assert.Equal(t, result, test.result)
	}
}

func TestIsPSQLService(t *testing.T) {
	tests := []isServiceTest{
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
	for _, test := range tests {
		serviceInstance := ServiceInstance{
			Service: test.serviceName,
			Plan:    test.planName,
		}
		result := serviceInstance.IsPSQLService()
		assert.Equal(t, result, test.result)
	}
}

type fetchServiceInstanceTest struct {
	serviceModel            plugin_models.GetService_Model
	getServiceError         error
	serviceName             string
	expectedServiceInstance ServiceInstance
	expectedError           error
}

func TestFetchServiceInstance(t *testing.T) {
	tests := []fetchServiceInstanceTest{
		{
			serviceModel: plugin_models.GetService_Model{
				Guid: "test-guid",
				Name: "something-else",
				ServicePlan: plugin_models.GetService_ServicePlan{
					Name: "shared-plan",
				},
				ServiceOffering: plugin_models.GetService_ServiceFields{
					Name: "aws-rds",
				},
			},
			getServiceError: nil,
			serviceName:     "my-test-db",
			expectedServiceInstance: ServiceInstance{
				GUID:    "test-guid",
				Service: "aws-rds",
				Name:    "my-test-db",
				Plan:    "shared-plan",
			},
		},
	}
	for _, test := range tests {
		fakeCliConnection := &pluginfakes.FakeCliConnection{}
		fakeCliConnection.GetServiceReturns(test.serviceModel, test.getServiceError)
		serviceInstance, err := FetchServiceInstance(fakeCliConnection, test.serviceName)
		assert.Equal(t, serviceInstance, test.expectedServiceInstance)
		assert.Equal(t, err, test.expectedError)
	}
}
