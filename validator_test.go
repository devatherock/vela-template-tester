package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(test *testing.T) {
	validationRequest := ValidationRequest{}

	content, error := ioutil.ReadFile("templates/input_template.yml")
	if error != nil {
		assert.Fail(test, error.Error())
	}
	validationRequest.Template = string(content)

	parameters := map[string]interface{}{
		"notification_branch": "develop",
		"notification_event":  "push",
	}
	validationRequest.Parameters = parameters

	validationResponse := validate(validationRequest)
	assert.Equal(test, "template is a valid yaml", validationResponse.Message)
}
