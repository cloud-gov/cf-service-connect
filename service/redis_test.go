package service

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/18F/cf-service-connect/models"
	"github.com/stretchr/testify/assert"
)

type redisMatchTest struct {
	serviceName string
	planName    string
	expected    bool
}

type mockCredentials struct {
	mockPassword string
}

func (m mockCredentials) GetPassword() string {
	return m.mockPassword
}

func (m mockCredentials) GetDBName() string {
	return ""
}

func (m mockCredentials) GetUsername() string {
	return ""
}

func (m mockCredentials) GetHost() string {
	return ""
}

func (m mockCredentials) GetPort() string {
	return ""
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

func TestGetRedisLaunchFlags(t *testing.T) {
	mockCredentialsClient := &mockCredentials{
		mockPassword: "fake-password",
	}
	flags := getRedisLaunchFlags(53839, mockCredentialsClient)
	expectedFlags := []string{
		"--tls",
		"-p", strconv.Itoa(53839),
		"-a", "fake-password",
	}
	if !reflect.DeepEqual(flags, expectedFlags) {
		t.Fatalf("expected: %s, got: %s", expectedFlags, flags)
	}
}
