// +build integration

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainWithCommandLineParameters(test *testing.T) {
	cases := []struct {
		parameters       map[string]string
		expectedExitCode int
		expectedOutput   string
	}{
		{
			map[string]string{
				"--templates": `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`,
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
		{
			map[string]string{
				"-ts": `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`,
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
		{
			map[string]string{
				"--input-file":      "templates/input_template.yml",
				"--variables":       `{"notification_branch":"develop","notification_event":"push"}`,
				"--expected-output": "templates/output_template.yml",
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
		{
			map[string]string{
				"-tf": "templates/input_template.yml",
				"-v":  `{"notification_branch":"develop","notification_event":"push"}`,
				"-o":  "templates/output_template.yml",
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
	}

	for _, data := range cases {
		arguments := []string{}
		for key, value := range data.parameters {
			arguments = append(arguments, key)
			arguments = append(arguments, value)
		}

		exitCode, output := executeCommand(exec.Command("./docker/velatemplatetesterplugin", arguments...))
		assert.Contains(test, output, data.expectedOutput)
		assert.Equal(test, data.expectedExitCode, exitCode)
	}
}

func TestMainWithEnvVariables(test *testing.T) {
	cases := []struct {
		parameters       map[string]string
		expectedExitCode int
		expectedOutput   string
	}{
		{
			map[string]string{
				"TEMPLATES": `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`,
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
		{
			map[string]string{
				"PARAMETER_TEMPLATES": `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`,
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
		{
			map[string]string{
				"INPUT_FILE":      "templates/input_template.yml",
				"VARIABLES":       `{"notification_branch":"develop","notification_event":"push"}`,
				"EXPECTED_OUTPUT": "templates/output_template.yml",
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE":      "templates/input_template.yml",
				"PARAMETER_VARIABLES":       `{"notification_branch":"develop","notification_event":"push"}`,
				"PARAMETER_EXPECTED_OUTPUT": "templates/output_template.yml",
			},
			0,
			"Template 'templates/input_template.yml' is valid.",
		},
	}

	for _, data := range cases {
		for key, value := range data.parameters {
			setEnvironmentVariable(test, key, value)
		}

		exitCode, output := executeCommand(exec.Command("./docker/velatemplatetesterplugin"))
		assert.Contains(test, output, data.expectedOutput)
		assert.Equal(test, data.expectedExitCode, exitCode)

		for key := range data.parameters {
			os.Unsetenv(key)
		}
	}
}
