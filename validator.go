package main

import (
    "bytes"
    "regexp"
    "strings"
    "text/template"

    "github.com/Masterminds/sprig"
    "gopkg.in/yaml.v2"
    log "github.com/sirupsen/logrus"
)

type ValidationResponse struct {
    Message string
    Error string `yaml:",omitempty"`
    Template string `yaml:",omitempty"`
}

type ValidationRequest struct {
    Parameters map[string]interface{}
    Template string
}

func validate(validationRequest ValidationRequest) ValidationResponse {
    // Process template
    buffer := new(bytes.Buffer)
    parsedTemplate, err := template.New("test").Funcs(sprig.TxtFuncMap()).Parse(validationRequest.Template)
    err = parsedTemplate.Execute(buffer, validationRequest.Parameters)

    validationResponse := ValidationResponse{}
    if err != nil {
        validationResponse.Error = err.Error()
        validationResponse.Message = "Invalid template"
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
