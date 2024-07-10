package main

import (
	"strings"
	"testing"

	"github.com/cloud-gov/cf-service-connect/connector"
	"github.com/stretchr/testify/assert"
)

type parseOptionsTest struct {
	args            string
	expectError     bool
	expectedOptions connector.Options
}

func TestParseOptions(t *testing.T) {
	tests := []parseOptionsTest{
		{
			"connect-to-service app service",
			false,
			connector.Options{
				"app",
				"service",
				true,
			},
		},
		{
			"connect-to-service -no-client app service",
			false,
			connector.Options{
				"app",
				"service",
				false,
			},
		},
		{
			"connect-to-service foo bar baz",
			true,
			connector.Options{},
		},
	}

	plugin := ServiceConnectPlugin{}
	for _, test := range tests {
		args := strings.Split(test.args, " ")
		opts, err := plugin.parseOptions(args)
		if test.expectError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, opts, test.expectedOptions)
		}
	}
}
