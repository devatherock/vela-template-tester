//go:build test
// +build test

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devatherock/vela-template-tester/pkg/validator"
	"github.com/devatherock/vela-template-tester/test/helper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestExpandTemplate(test *testing.T) {
	cases := []struct {
		inputFile    string
		parameters   map[string]interface{}
		outputFile   string
		templateType string
	}{
		{
			"test/testdata/input_template.yml",
			map[string]interface{}{
				"notification_branch": "develop",
				"notification_event":  "push",
			},
			"test/testdata/output_template.yml",
			"",
		},
		{
			"test/testdata/input_starlark_template.py",
			map[string]interface{}{
				"image": "go:1.14",
			},
			"test/testdata/output_starlark_template.yml",
			"starlark",
		},
	}

	for _, data := range cases {
		validationRequest := validator.ValidationRequest{}
		input, _ := ioutil.ReadFile(helper.AbsolutePath(data.inputFile))
		validationRequest.Template = string(input)
		validationRequest.Type = data.templateType
		validationRequest.Parameters = data.parameters

		yamlStr, _ := yaml.Marshal(&validationRequest)
		request, _ := http.NewRequest("POST", "/api/expandTemplate", bytes.NewBuffer(yamlStr))

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(expandTemplate)
		handler.ServeHTTP(response, request)

		assert.Equal(test, 200, response.Code)

		validationResponse := validator.ValidationResponse{}
		yaml.Unmarshal(response.Body.Bytes(), &validationResponse)
		assert.Equal(test, "template is a valid yaml", validationResponse.Message)
		assert.Equal(test, "", validationResponse.Error)

		expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath(data.outputFile))
		expectedOutputMap := make(map[interface{}]interface{})
		yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

		processedTemplateMap := make(map[interface{}]interface{})
		yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)

		assert.Equal(test, expectedOutputMap, processedTemplateMap)
	}
}

func TestExpandTemplateError(test *testing.T) {
	cases := []struct {
		inputFile     string
		parameters    map[string]interface{}
		message       string
		expectedError string
	}{
		{
			"test/testdata/input_parse_error_template.yml",
			map[string]interface{}{
				"notification_branch": "develop",
			},
			"Invalid template",
			"Unable to parse template",
		},
		{
			"test/testdata/input_invalid_template.yml",
			map[string]interface{}{
				"notification_branch": "develop",
			},
			"template is not a valid yaml",
			"yaml: line 4: did not find expected ',' or ']'",
		},
	}

	for _, data := range cases {
		validationRequest := map[string]interface{}{}
		input, _ := ioutil.ReadFile(helper.AbsolutePath(data.inputFile))
		validationRequest["template"] = string(input)
		validationRequest["parameters"] = data.parameters

		yamlStr, _ := yaml.Marshal(&validationRequest)
		request, _ := http.NewRequest("POST", "/api/expandTemplate", bytes.NewBuffer(yamlStr))

		response := httptest.NewRecorder()
		handler := http.HandlerFunc(expandTemplate)
		handler.ServeHTTP(response, request)

		assert.Equal(test, 200, response.Code)

		validationResponse := map[string]string{}
		yaml.Unmarshal(response.Body.Bytes(), &validationResponse)
		assert.Equal(test, data.message, validationResponse["message"])
		assert.Equal(test, data.expectedError, validationResponse["error"])
	}
}

func TestExpandTemplateListParams(test *testing.T) {
	validationRequest := validator.ValidationRequest{}
	input, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/list_parameters_template.yml"))
	validationRequest.Template = string(input)

	parameters := []interface{}{
		map[string]string{
			"title":   "Hello",
			"content": "World",
		},
		map[string]string{
			"title":   "Hi",
			"content": "There",
		},
	}
	validationRequest.Parameters = parameters

	yamlStr, _ := yaml.Marshal(&validationRequest)
	request, _ := http.NewRequest("POST", "/api/expandTemplate", bytes.NewBuffer(yamlStr))

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(expandTemplate)
	handler.ServeHTTP(response, request)

	assert.Equal(test, 200, response.Code)

	validationResponse := validator.ValidationResponse{}
	yaml.Unmarshal(response.Body.Bytes(), &validationResponse)
	assert.Equal(test, "template is a valid yaml", validationResponse.Message)
	assert.Equal(test, "", validationResponse.Error)

	expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/list_parameters_output.yml"))
	expectedOutputMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

	processedTemplateMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)

	assert.Equal(test, expectedOutputMap, processedTemplateMap)
}

func TestCheckHealth(test *testing.T) {
	request, _ := http.NewRequest("GET", "/api/health", nil)

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(checkHealth)
	handler.ServeHTTP(response, request)

	assert.Equal(test, 200, response.Code)
	assert.Equal(test, "UP", response.Body.String())
}

func TestLookupPortEnvVariablePresent(test *testing.T) {
	helper.SetEnvironmentVariable(test, "PORT", "8081")

	assert.Equal(test, "8081", lookupPort())
}

func TestLookupPortEnvVariableAbsent(test *testing.T) {
	assert.Equal(test, "8080", lookupPort())
}
