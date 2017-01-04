package models

import (
	"testing"

	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/stretchr/testify/assert"
)

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
