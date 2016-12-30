package models

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"reflect"
	"testing"
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
		if result != test.result {
			t.Errorf("Expected result %v. Real result %v. Data: Service Name '%s' Plan Name '%s'",
				test.result, result, test.serviceName, test.planName)
		}
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
		if result != test.result {
			t.Errorf("Expected result %v. Real result %v. Data: Service Name '%s' Plan Name '%s'",
				test.result, result, test.serviceName, test.planName)
		}
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
		if !reflect.DeepEqual(serviceInstance, test.expectedServiceInstance) {
			t.Errorf("Failed Service Instance Equality Check. Expected %+v. Actual %+v\n", test.expectedServiceInstance, serviceInstance)
		}
		if err != test.expectedError {
			t.Errorf("Failed Returned Error Check. Expected %+v. Actual %+v\n", test.expectedError, err)
		}
	}
}
