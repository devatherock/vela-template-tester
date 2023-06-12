//go:build integration
// +build integration

package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/devatherock/vela-template-tester/test/helper"
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
				"--templates": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
				),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{
				"-ts": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
				),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{
				"--input-file":      helper.AbsolutePath("test/testdata/input_template.yml"),
				"--variables":       `{"notification_branch":"develop","notification_event":"push"}`,
				"--expected-output": helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{
				"-tf": helper.AbsolutePath("test/testdata/input_template.yml"),
				"-v":  `{"notification_branch":"develop","notification_event":"push"}`,
				"-o":  helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
	}

	for _, data := range cases {
		arguments := []string{}
		for key, value := range data.parameters {
			arguments = append(arguments, key)
			arguments = append(arguments, value)
		}

		exitCode, output := helper.ExecuteCommand(exec.Command(helper.AbsolutePath("bin/plugin"), arguments...))
		fmt.Print(output)

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
				"TEMPLATES": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
				),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{
				"PARAMETER_TEMPLATES": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
				),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{
				"INPUT_FILE":      helper.AbsolutePath("test/testdata/input_template.yml"),
				"VARIABLES":       `{"notification_branch":"develop","notification_event":"push"}`,
				"EXPECTED_OUTPUT": helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE":      helper.AbsolutePath("test/testdata/input_template.yml"),
				"PARAMETER_VARIABLES":       `{"notification_branch":"develop","notification_event":"push"}`,
				"PARAMETER_EXPECTED_OUTPUT": helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			0,
			fmt.Sprintf("Template '%s' is valid.", helper.AbsolutePath("test/testdata/input_template.yml")),
		},
		{
			map[string]string{},
			0,
			"No template specified",
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE": helper.AbsolutePath("test/testdata/input_parse_error_template.yml"),
				"PARAMETER_VARIABLES":  `{"notification_branch":"develop"}`,
			},
			1,
			fmt.Sprintf(
				"Template '%s' is invalid. Error: Unable to parse template",
				helper.AbsolutePath("test/testdata/input_parse_error_template.yml"),
			),
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE": helper.AbsolutePath("test/testdata/input_invalid_template.yml"),
				"PARAMETER_VARIABLES":  `{"notification_branch":"develop"}`,
			},
			1,
			fmt.Sprintf(
				"Template '%s' is invalid. Error: yaml: line 4: did not find expected ',' or ']'",
				helper.AbsolutePath("test/testdata/input_invalid_template.yml"),
			),
		},
		{
			map[string]string{
				"PARAMETER_INPUT_FILE":      helper.AbsolutePath("test/testdata/input_template.yml"),
				"PARAMETER_EXPECTED_OUTPUT": helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			1,
			fmt.Sprintf(
				"Template '%s' is valid, but did not match expected output",
				helper.AbsolutePath("test/testdata/input_template.yml"),
			),
		},
	}

	for _, data := range cases {
		for key, value := range data.parameters {
			helper.SetEnvironmentVariable(test, key, value)
		}

		exitCode, output := helper.ExecuteCommand(exec.Command(helper.AbsolutePath("bin/plugin")))
		fmt.Print(output)

		assert.Contains(test, output, data.expectedOutput)
		assert.Equal(test, data.expectedExitCode, exitCode)

		for key := range data.parameters {
			os.Unsetenv(key)
		}
	}
}
