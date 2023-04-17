package service

import (
	"testing"

	"github.com/bruth/assert"
)

type launchCmdStringTest struct {
	LaunchCmd LaunchCmd
	Expected  string
}

func TestLaunchCmdString(t *testing.T) {
	tests := []launchCmdStringTest{
		launchCmdStringTest{
			LaunchCmd: LaunchCmd{
				Cmd: "foo",
			},
			Expected: "foo",
		},
		launchCmdStringTest{
			LaunchCmd: LaunchCmd{
				Envs: map[string]string{
					"a": "1",
					"b": "2",
				},
				Cmd: "foo",
			},
			Expected: "a=1 b=2 foo",
		},
		launchCmdStringTest{
			LaunchCmd: LaunchCmd{
				Cmd:  "foo",
				Args: []string{"bar", "baz"},
			},
			Expected: "foo bar baz",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.LaunchCmd.String(), test.Expected)
	}
}
