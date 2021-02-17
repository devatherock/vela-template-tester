// +build !integration

package main

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ValidationResponse struct {
	Message  string
	Error    string `yaml:",omitempty"`
	Template string `yaml:",omitempty"`
}

type ValidationRequest struct {
	Parameters map[string]interface{}
	Template   string
}

func validate(validationRequest ValidationRequest) (validationResponse ValidationResponse) {
	validationResponse = ValidationResponse{}

	// Error response in case of a panic
	validationResponse.Message = "Invalid template"
	defer handlePanic()

	// Process template
	buffer := new(bytes.Buffer)
	parsedTemplate, err := template.New("test").Funcs(sprig.TxtFuncMap()).Parse(validationRequest.Template)
	err = parsedTemplate.Execute(buffer, validationRequest.Parameters)

	if err != nil {
		validationResponse.Error = err.Error()
	} else {
		processedTemplate := make(map[interface{}]interface{})

		// To prevent yaml from being output in flow style due to trailing spaces
		regex := regexp.MustCompile(`[^\S\r\n]+\n`)
		validationResponse.Template = strings.TrimSpace(regex.ReplaceAllString(buffer.String(), "\n"))

		err = yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplate)
		if err != nil {
			validationResponse.Error = err.Error()
			validationResponse.Message = "template is not a valid yaml"
		} else {
			validationResponse.Message = "template is a valid yaml"
		}
		log.Debug("Output template: \n", buffer.String())
	}

	return validationResponse
}

func handlePanic() {
	if error := recover(); error != nil {
		log.Error("Recovering from panic: ", error)
	}
}
