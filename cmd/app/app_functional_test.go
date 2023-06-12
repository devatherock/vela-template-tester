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
	validationRequest := map[string]interface{}{}
	input, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/input_template.yml"))
	validationRequest["template"] = string(input)

	parameters := map[string]interface{}{
		"notification_branch": "develop",
		"notification_event":  "push",
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

	expectedOutput, _ := ioutil.ReadFile(helper.AbsolutePath("test/testdata/output_template.yml"))
	expectedOutputMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)

	processedTemplateMap := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(validationResponse["template"]), &processedTemplateMap)

	assert.Equal(test, expectedOutputMap, processedTemplateMap)
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
