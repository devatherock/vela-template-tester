//go:build functional
// +build functional

package main

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/devatherock/vela-template-tester/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestMainWithEnvVariables(test *testing.T) {
	cases := []struct {
		parameters       map[string]string
		expectedExitCode int
		expectedOutput   string
	}{
		{
			map[string]string{
				"TEMPLATES": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					"test/testdata/input_template.yml",
					"test/testdata/output_template.yml",
				),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", "test/testdata/input_template.yml"),
		},
		{
			map[string]string{
				"PARAMETER_TEMPLATES": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					"test/testdata/input_template.yml",
					"test/testdata/output_template.yml",
				),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", "test/testdata/input_template.yml"),
		},
		{
			map[string]string{
				"INPUT_FILE":      "test/testdata/input_template.yml",
				"VARIABLES":       `{"notification_branch":"develop","notification_event":"push"}`,
				"EXPECTED_OUTPUT": "test/testdata/output_template.yml",
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", "test/testdata/input_template.yml"),
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE":      "test/testdata/input_template.yml",
				"PARAMETER_VARIABLES":       `{"notification_branch":"develop","notification_event":"push"}`,
				"PARAMETER_EXPECTED_OUTPUT": "test/testdata/output_template.yml",
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", "test/testdata/input_template.yml"),
		},
		{
			map[string]string{},
			0,
			"No template specified",
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE": "test/testdata/input_parse_error_template.yml",
				"PARAMETER_VARIABLES":  `{"notification_branch":"develop"}`,
			},
			1,
			fmt.Sprintf(
				"Template '%s' is invalid. Error: Unable to parse template",
				"test/testdata/input_parse_error_template.yml",
			),
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE": "test/testdata/input_invalid_template.yml",
				"PARAMETER_VARIABLES":  `{"notification_branch":"develop"}`,
			},
			1,
			fmt.Sprintf(
				"Template '%s' is invalid. Error: yaml: line 4: did not find expected ',' or ']'",
				"test/testdata/input_invalid_template.yml",
			),
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE":      "test/testdata/input_template.yml",
				"PARAMETER_EXPECTED_OUTPUT": "test/testdata/output_template.yml",
			},
			1,
			fmt.Sprintf(
				"Template '%s' is valid, but did not match expected output",
				"test/testdata/input_template.yml",
			),
		},
	}

	for _, data := range cases {
		arguments := []string{"run", "--rm", "-v",
			helper.GetProjectRoot() + ":/work",
			"-w", "/work",
		}

		for key, value := range data.parameters {
			arguments = append(arguments, "-e")
			arguments = append(arguments, fmt.Sprintf("%s=%s", key, value))
		}
		arguments = append(arguments, "devatherock/vela-template-tester:latest")

		exitCode, output := helper.ExecuteCommand(exec.Command("docker", arguments...))
		fmt.Print(output)

		assert.Contains(test, output, data.expectedOutput)
		assert.Equal(test, data.expectedExitCode, exitCode)
	}
}
