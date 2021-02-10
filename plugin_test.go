package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyOutput(test *testing.T) {
	request := PluginValidationRequest{}
	validationResponse := ValidationResponse{}

	assert.True(test, verifyOutput(request, validationResponse))
}
