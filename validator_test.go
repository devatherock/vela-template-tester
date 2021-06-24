// +build !integration

package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestValidateSuccess(test *testing.T) {
	validationRequest := ValidationRequest{}

	input, _ := ioutil.ReadFile("templates/input_template.yml")
	validationRequest.Template = string(input)

	parameters := map[string]interface{}{
		"notification_branch": "develop",
		"notification_event":  "push",
	}
	validationRequest.Parameters = parameters

	validationResponse := validate(validationRequest)
	assert.Equal(test, "template is a valid yaml", validationResponse.Message)
	assert.Equal(test, "", validationResponse.Error)

	expectedOutput, _ := ioutil.ReadFile("templates/output_template.yml")
	expectedOutputMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

	processedTemplateMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)

	assert.Equal(test, expectedOutputMap, processedTemplateMap)
}

func TestValidateVelaFunctionSuccess(test *testing.T) {
	validationRequest := ValidationRequest{}

	input, _ := ioutil.ReadFile("templates/input_vela_function_template.yml")
	validationRequest.Template = string(input)

	validationResponse := validate(validationRequest)
	assert.Equal(test, "template is a valid yaml", validationResponse.Message)
	assert.Equal(test, "", validationResponse.Error)

	expectedOutput, _ := ioutil.ReadFile("templates/output_vela_function_template.yml")
	expectedOutputMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

	processedTemplateMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)

	assert.Equal(test, expectedOutputMap, processedTemplateMap)
}

func TestValidateVelaFunctionFailure(test *testing.T) {
	validationRequest := ValidationRequest{}

	input, _ := ioutil.ReadFile("templates/input_vela_fn_empty_variable_template.yml")
	validationRequest.Template = string(input)

	validationResponse := validate(validationRequest)
	assert.Equal(test, "Invalid template", validationResponse.Message)
	assert.Contains(test, validationResponse.Error, "Environment variable name cannot be empty in 'vela' function")
}

func TestValidateParseError(test *testing.T) {
	validationRequest := ValidationRequest{}

	input, _ := ioutil.ReadFile("templates/input_parse_error_template.yml")
	validationRequest.Template = string(input)

	parameters := map[string]interface{}{
		"notification_branch": "develop",
	}
	validationRequest.Parameters = parameters

	validationResponse := validate(validationRequest)
	assert.Equal(test, "Invalid template", validationResponse.Message)
	assert.Equal(test, "Unable to parse template", validationResponse.Error)
	assert.Equal(test, "", validationResponse.Template)
}

func TestValidateInvalidTemplate(test *testing.T) {
	validationRequest := ValidationRequest{}

	input, _ := ioutil.ReadFile("templates/input_invalid_template.yml")
	validationRequest.Template = string(input)

	parameters := map[string]interface{}{
		"notification_branch": "develop",
	}
	validationRequest.Parameters = parameters

	validationResponse := validate(validationRequest)
	assert.Equal(test, "template is not a valid yaml", validationResponse.Message)
	assert.Equal(test, "yaml: line 4: did not find expected ',' or ']'", validationResponse.Error)

	expectedOutput, _ := ioutil.ReadFile("templates/output_invalid_template.txt")
	assert.Equal(test, string(expectedOutput), validationResponse.Template)
}
