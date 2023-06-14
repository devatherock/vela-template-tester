//go:build api && !test
// +build api,!test

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/devatherock/vela-template-tester/test/helper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const baseUrl string = "http://localhost:8082"

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
		validationRequest := map[string]interface{}{}
		input, _ := ioutil.ReadFile(helper.AbsolutePath(data.inputFile))
		validationRequest["template"] = string(input)
		validationRequest["parameters"] = data.parameters
		validationRequest["type"] = data.templateType

		yamlStr, _ := yaml.Marshal(&validationRequest)
		request, _ := http.NewRequest("POST", baseUrl+"/api/expandTemplate", bytes.NewBuffer(yamlStr))

		client := &http.Client{}
		response, err := client.Do(request)
		defer response.Body.Close()

		assert.Nil(test, err)
		assert.Equal(test, 200, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		assert.Nil(test, err)

		validationResponse := map[string]string{}
		yaml.Unmarshal(responseBody, &validationResponse)
		assert.Equal(test, "template is a valid yaml", validationResponse["message"])
		assert.Equal(test, "", validationResponse["error"])

		expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath(data.outputFile))
		expectedOutputMap := make(map[interface{}]interface{})
		yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

		processedTemplateMap := make(map[interface{}]interface{})
		yaml.Unmarshal([]byte(validationResponse["template"]), &processedTemplateMap)

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
		request, _ := http.NewRequest("POST", baseUrl+"/api/expandTemplate", bytes.NewBuffer(yamlStr))

		client := &http.Client{}
		response, err := client.Do(request)
		defer response.Body.Close()

		assert.Nil(test, err)
		assert.Equal(test, 200, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		assert.Nil(test, err)

		validationResponse := map[string]string{}
		yaml.Unmarshal(responseBody, &validationResponse)
		assert.Equal(test, data.message, validationResponse["message"])
		assert.Equal(test, data.expectedError, validationResponse["error"])
	}
}

func TestExpandTemplateListParams(test *testing.T) {
	validationRequest := map[string]interface{}{}
	input, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/list_parameters_template.yml"))
	validationRequest["template"] = string(input)

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
	validationRequest["parameters"] = parameters

	yamlStr, _ := yaml.Marshal(&validationRequest)
	request, _ := http.NewRequest("POST", baseUrl+"/api/expandTemplate", bytes.NewBuffer(yamlStr))

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()

	assert.Nil(test, err)
	assert.Equal(test, 200, response.StatusCode)

	responseBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(test, err)

	validationResponse := map[string]string{}
	yaml.Unmarshal(responseBody, &validationResponse)
	assert.Equal(test, "template is a valid yaml", validationResponse["message"])
	assert.Equal(test, "", validationResponse["error"])

	expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/list_parameters_output.yml"))
	expectedOutputMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

	processedTemplateMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(validationResponse["template"]), &processedTemplateMap)

	assert.Equal(test, expectedOutputMap, processedTemplateMap)
}

func TestCheckHealth(test *testing.T) {
	request, _ := http.NewRequest("GET", baseUrl+"/api/health", nil)

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()

	assert.Nil(test, err)
	assert.Equal(test, 200, response.StatusCode)

	responseBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(test, err)
	assert.Equal(test, "UP", string(responseBody))
}
