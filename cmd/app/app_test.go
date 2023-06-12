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
	validationRequest := validator.ValidationRequest{}
	input, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/input_template.yml"))
	validationRequest.Template = string(input)

	parameters := map[string]interface{}{
		"notification_branch": "develop",
		"notification_event":  "push",
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

	expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/output_template.yml"))
	expectedOutputMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

	processedTemplateMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)

	assert.Equal(test, expectedOutputMap, processedTemplateMap)
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
