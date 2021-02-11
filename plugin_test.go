// +build !api

package main

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestReadInputParameters(test *testing.T) {
	cases := []struct {
		parameters map[string]string
		expected   []PluginValidationRequest
	}{
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
