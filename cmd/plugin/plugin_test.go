//go:build test
// +build test

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/devatherock/vela-template-tester/pkg/validator"
	"github.com/devatherock/vela-template-tester/test/helper"
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
		helper.SetEnvironmentVariable(test, data.inputFileEnv, helper.AbsolutePath("test/testdata/input_template.yml"))
		helper.SetEnvironmentVariable(test, data.variablesEnv, `{"notification_branch":"develop","notification_event":"push"}`)
		helper.SetEnvironmentVariable(test, data.expectedOutputEnv, helper.AbsolutePath("test/testdata/output_template.yml"))

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
		helper.SetEnvironmentVariable(test, data, fmt.Sprintf(
			`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
			helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
		))
		runApp([]string{"-x", "dummy"})

		os.Unsetenv(data)
	}
}

func TestRunAppWithIndividualParameters(test *testing.T) {
	exitCode := captureExitCode(test)

	cases := []struct {
		inputFileParam           string
		variablesParam           string
		expectedOutputParam      string
		templateTypeParam        string
		inputFileParamValue      string
		variablesParamValue      string
		expectedOutputParamValue string
		templateTypeParamValue   string
	}{
		{
			"--input-file",
			"--variables",
			"--expected-output",
			"--template-type",
			"test/testdata/input_template.yml",
			`{"notification_branch":"develop","notification_event":"push"}`,
			"test/testdata/output_template.yml",
			"",
		},
		{
			"-tf",
			"-v",
			"-o",
			"-tt",
			"test/testdata/input_starlark_template.py",
			`{"image":"go:1.14"}`,
			"test/testdata/output_starlark_template.yml",
			"starlark",
		},
	}

	for _, data := range cases {
		arguments := []string{
			data.inputFileParam, helper.AbsolutePath(data.inputFileParamValue),
			data.variablesParam, data.variablesParamValue,
			data.expectedOutputParam, helper.AbsolutePath(data.expectedOutputParamValue),
			data.templateTypeParam, data.templateTypeParamValue,
		}

		runApp(arguments)

		assert.Equal(test, 0, exitCode[0])
	}
}

func TestRunAppWithTemplatesParameter(test *testing.T) {
	exitCode := captureExitCode(test)

	cases := []struct {
		templatesParam      string
		templatesParamValue string
	}{
		{
			"--templates",
			fmt.Sprintf(
				`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
				helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
			),
		},
		{
			"-ts",
			fmt.Sprintf(
				`[{"input_file":"%s","variables":{"image":"go:1.14"},"expected_output":"%s","template_type":"starlark"}]`,
				helper.AbsolutePath("test/testdata/input_starlark_template.py"),
				helper.AbsolutePath("test/testdata/output_starlark_template.yml"),
			),
		},
	}

	for _, data := range cases {
		runApp([]string{data.templatesParam, data.templatesParamValue})

		assert.Equal(test, 0, exitCode[0])
	}
}

func TestRun(test *testing.T) {
	exitCode := captureExitCode(test)

	cases := []struct {
		parameters       map[string]string
		expected         error
		expectedExitCode int
	}{
		{
			map[string]string{
				"input-file": helper.AbsolutePath("test/testdata/input_template.yml"),
			},
			nil,
			-1,
		},
		{
			map[string]string{
				"input-file": helper.AbsolutePath("test/testdata/input_invalid_template.yml"),
				"variables":  `{"notification_branch":"develop"}`,
			},
			fmt.Errorf(
				"Template '%s' is invalid. Error: yaml: line 4: did not find expected ',' or ']'",
				helper.AbsolutePath("test/testdata/input_invalid_template.yml"),
			),
			1,
		},
		{
			map[string]string{
				"input-file":      helper.AbsolutePath("test/testdata/input_template.yml"),
				"expected-output": helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			fmt.Errorf(
				"Template '%s' is valid, but did not match expected output",
				helper.AbsolutePath("test/testdata/input_template.yml"),
			),
			1,
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
		assert.Equal(test, data.expectedExitCode, exitCode[0])
	}
}

func TestReadInputParameters(test *testing.T) {
	captureExitCode(test)

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
				"input-file": helper.AbsolutePath("test/testdata/input_template.yml"),
			},
			[]PluginValidationRequest{
				{
					InputFile: helper.AbsolutePath("test/testdata/input_template.yml"),
				},
			},
		},
		{
			map[string]string{
				"input-file": helper.AbsolutePath("test/testdata/input_template.yml"),
				"variables":  `{"notification_branch":"develop","notification_event":"push"}`,
			},
			[]PluginValidationRequest{
				{
					InputFile: helper.AbsolutePath("test/testdata/input_template.yml"),
					Variables: map[string]interface{}{
						"notification_branch": "develop",
						"notification_event":  "push",
					},
				},
			},
		},
		{
			map[string]string{
				"input-file":      helper.AbsolutePath("test/testdata/input_template.yml"),
				"variables":       `{"notification_branch":"develop","notification_event":"push"}`,
				"expected-output": helper.AbsolutePath("test/testdata/output_template.yml"),
			},
			[]PluginValidationRequest{
				{
					InputFile: helper.AbsolutePath("test/testdata/input_template.yml"),
					Variables: map[string]interface{}{
						"notification_branch": "develop",
						"notification_event":  "push",
					},
					ExpectedOutput: helper.AbsolutePath("test/testdata/output_template.yml"),
				},
			},
		},
		{
			map[string]string{
				"templates": fmt.Sprintf(
					`[{"input_file":"%s","variables":{"notification_branch":"develop","notification_event":"push"},"expected_output":"%s"}]`,
					helper.AbsolutePath("test/testdata/input_template.yml"), helper.AbsolutePath("test/testdata/output_template.yml"),
				),
			},
			[]PluginValidationRequest{
				{
					InputFile: helper.AbsolutePath("test/testdata/input_template.yml"),
					Variables: map[string]interface{}{
						"notification_branch": "develop",
						"notification_event":  "push",
					},
					ExpectedOutput: helper.AbsolutePath("test/testdata/output_template.yml"),
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
	validationResponse := validator.ValidationResponse{}

	assert.True(test, verifyOutput(request, validationResponse))
}

func TestVerifyOutputSuccess(test *testing.T) {
	request := PluginValidationRequest{}
	request.ExpectedOutput = helper.AbsolutePath("test/testdata/output_template.yml")

	validationResponse := validator.ValidationResponse{}
	expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/output_template.yml"))
	validationResponse.Template = string(expectedOutput)

	assert.True(test, verifyOutput(request, validationResponse))
}

func TestVerifyOutputFailure(test *testing.T) {
	request := PluginValidationRequest{}
	request.ExpectedOutput = helper.AbsolutePath("test/testdata/output_template.yml")

	validationResponse := validator.ValidationResponse{}
	validationResponse.Template = "foo: bar"

	assert.False(test, verifyOutput(request, validationResponse))
}

// Overrides exit code function for tests
func captureExitCode(test *testing.T) []int {
	originalExitFunction := exit
	exitCode := []int{-1}
	exit = func(code int) { exitCode[0] = code }

	test.Cleanup(func() {
		exit = originalExitFunction
	})

	return exitCode
}
