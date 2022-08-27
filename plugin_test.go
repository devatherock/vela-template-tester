//go:build !api && !integration
// +build !api,!integration

package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestRunAppWithIndividualEnvVariables(test *testing.T) {
	cases := []struct {
		inputFileEnv      string
		variablesEnv      string
		expectedOutputEnv string
	}{
		{
			"PARAMETER_INPUT_FILE",
			"PARAMETER_VARIABLES",
			"PARAMETER_EXPECTED_OUTPUT",
		},
		{
			"INPUT_FILE",
			"VARIABLES",
			"EXPECTED_OUTPUT",
		},
	}

	for _, data := range cases {
		setEnvironmentVariable(test, data.inputFileEnv, "templates/input_template.yml")
		setEnvironmentVariable(test, data.variablesEnv, `{"notification_branch":"develop","notification_event":"push"}`)
		setEnvironmentVariable(test, data.expectedOutputEnv, "templates/output_template.yml")

		runApp([]string{"-x", "dummy"})

		os.Unsetenv(data.inputFileEnv)
		os.Unsetenv(data.variablesEnv)
		os.Unsetenv(data.expectedOutputEnv)
	}
}

func TestRunAppWithTemplatesEnvVariable(test *testing.T) {
	cases := []string{
		"PARAMETER_TEMPLATES",
		"TEMPLATES",
	}

	for _, data := range cases {
		setEnvironmentVariable(test, data, `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`)
		runApp([]string{"-x", "dummy"})

		os.Unsetenv(data)
	}
}

func TestRunAppWithIndividualParameters(test *testing.T) {
	cases := []struct {
		inputFileParam      string
		variablesParam      string
		expectedOutputParam string
	}{
		{
			"--input-file",
			"--variables",
			"--expected-output",
		},
		{
			"-tf",
			"-v",
			"-o",
		},
	}

	for _, data := range cases {
		arguments := []string{
			data.inputFileParam, "templates/input_template.yml",
			data.variablesParam, `{"notification_branch":"develop","notification_event":"push"}`,
			data.expectedOutputParam, "templates/output_template.yml",
		}

		runApp(arguments)
	}
}

func TestRunAppWithTemplatesParameter(test *testing.T) {
	cases := []string{
		"--templates",
		"-ts",
	}

	for _, data := range cases {
		runApp([]string{data, `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`})
	}
}

func TestRun(test *testing.T) {
	cases := []struct {
		parameters map[string]string
		expected   error
	}{
		{
			map[string]string{
				"input-file": "templates/input_template.yml",
			},
			nil,
		},
		{
			map[string]string{
				"input-file": "templates/input_invalid_template.yml",
				"variables":  `{"notification_branch":"develop"}`,
			},
			errors.New("Template 'templates/input_invalid_template.yml' is invalid. Error: yaml: line 4: did not find expected ',' or ']'"),
		},
		{
			map[string]string{
				"input-file":      "templates/input_template.yml",
				"expected-output": "templates/output_template.yml",
			},
			errors.New("Template 'templates/input_template.yml' is valid, but did not match expected output"),
		},
	}

	for _, data := range cases {
		set := flag.NewFlagSet("test", 0)
		for key, value := range data.parameters {
			set.String(key, value, "")
		}

		context := cli.NewContext(nil, set, nil)
		actual := run(context)

		assert.Equal(test, data.expected, actual)
	}
}

func TestReadInputParameters(test *testing.T) {
	cases := []struct {
		parameters map[string]string
		expected   []PluginValidationRequest
	}{
		{
			map[string]string{},
			[]PluginValidationRequest{},
		},
		{
			map[string]string{
				"input-file": "templates/input_template.yml",
			},
			[]PluginValidationRequest{
				{
					InputFile: "templates/input_template.yml",
				},
			},
		},
		{
			map[string]string{
				"input-file": "templates/input_template.yml",
				"variables":  `{"notification_branch":"develop","notification_event":"push"}`,
			},
			[]PluginValidationRequest{
				{
					InputFile: "templates/input_template.yml",
					Variables: map[string]interface{}{
						"notification_branch": "develop",
						"notification_event":  "push",
					},
				},
			},
		},
		{
			map[string]string{
				"input-file":      "templates/input_template.yml",
				"variables":       `{"notification_branch":"develop","notification_event":"push"}`,
				"expected-output": "templates/output_template.yml",
			},
			[]PluginValidationRequest{
				{
					InputFile: "templates/input_template.yml",
					Variables: map[string]interface{}{
						"notification_branch": "develop",
						"notification_event":  "push",
					},
					ExpectedOutput: "templates/output_template.yml",
				},
			},
		},
		{
			map[string]string{
				"templates": `[{"input_file":"templates/input_template.yml","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"templates/output_template.yml"}]`,
			},
			[]PluginValidationRequest{
				{
					InputFile: "templates/input_template.yml",
					Variables: map[string]interface{}{
						"notification_branch": "develop",
						"notification_event":  "push",
					},
					ExpectedOutput: "templates/output_template.yml",
				},
			},
		},
	}

	for _, data := range cases {
		set := flag.NewFlagSet("test", 0)
		for key, value := range data.parameters {
			set.String(key, value, "")
		}

		context := cli.NewContext(nil, set, nil)
		actual := readInputParameters(context)

		assert.Equal(test, data.expected, actual)
	}
}

func TestVerifyOutputDisabled(test *testing.T) {
	request := PluginValidationRequest{}
	validationResponse := ValidationResponse{}

	assert.True(test, verifyOutput(request, validationResponse))
}

func TestVerifyOutputSuccess(test *testing.T) {
	request := PluginValidationRequest{}
	request.ExpectedOutput = "templates/output_template.yml"

	validationResponse := ValidationResponse{}
	expectedOutput, _ := ioutil.ReadFile("templates/output_template.yml")
	validationResponse.Template = string(expectedOutput)

	assert.True(test, verifyOutput(request, validationResponse))
}

func TestVerifyOutputFailure(test *testing.T) {
	request := PluginValidationRequest{}
	request.ExpectedOutput = "templates/output_template.yml"

	validationResponse := ValidationResponse{}
	validationResponse.Template = "foo: bar"

	assert.False(test, verifyOutput(request, validationResponse))
}
