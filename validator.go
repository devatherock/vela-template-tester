// +build !integration

package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	log "github.com/sirupsen/logrus"
	"github.com/starlight-go/starlight/convert"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"
)

type ValidationResponse struct {
	Message  string
	Error    string `yaml:",omitempty"`
	Template string `yaml:",omitempty"`
}

type ValidationRequest struct {
	Parameters interface{}
	Template   string
	Type       string
}

func validate(validationRequest ValidationRequest) (validationResponse ValidationResponse) {
	validationResponse = ValidationResponse{}

	// Error response in case of a panic
	validationResponse.Message = "Invalid template"
	validationResponse.Error = "Unable to parse template"
	defer handlePanic()

	// Process template
	var outputTemplate string
	var err error
	if validationRequest.Type == "starlark" {
		outputTemplate, err = validateStarlarkTemplate(&validationRequest)
	} else {
		outputTemplate, err = validateGoTemplate(&validationRequest)
	}

	if err != nil {
		validationResponse.Error = err.Error()
	} else {
		processedTemplate := make(map[interface{}]interface{})

		// To prevent yaml from being output in flow style due to trailing spaces
		regex := regexp.MustCompile(`[^\S\r\n]+\n`)
		validationResponse.Template = strings.TrimSpace(regex.ReplaceAllString(outputTemplate, "\n"))

		err = yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplate)
		if err != nil {
			validationResponse.Error = err.Error()
			validationResponse.Message = "template is not a valid yaml"
		} else {
			validationResponse.Message = "template is a valid yaml"
			validationResponse.Error = ""
		}
		log.Debug("Output template: \n", outputTemplate)
	}

	return validationResponse
}

func validateGoTemplate(validationRequest *ValidationRequest) (string, error) {
	buffer := new(bytes.Buffer)
	parsedTemplate, err := template.New("test").Funcs(VelaFuncMap()).Funcs(sprig.TxtFuncMap()).Parse(validationRequest.Template)
	err = parsedTemplate.Execute(buffer, validationRequest.Parameters)

	outputTemplate := buffer.String()
	return outputTemplate, err
}

func validateStarlarkTemplate(validationRequest *ValidationRequest) (string, error) {
	context := map[string]interface{}{
		"ctx": map[string]interface{}{
			"vars": validationRequest.Parameters,
		},
		"Println": fmt.Println,
	}
	thread := &starlark.Thread{
		Load: nil,
	}
	dict, _ := MakeStarlarkStringDict(context)
	output, err := starlark.ExecFile(thread, "vela.template"+time.Now().String(), validationRequest.Template, dict)
	fmt.Println(output)

	if err != nil {
		return "", err
	}

	fmt.Println(reflect.TypeOf(output["output"]))
	steps := output["output"].(*starlark.Dict)
	stepsKey, _ := convert.ToValue("steps")
	stepsValue, _, _ := steps.Get(stepsKey)
	fmt.Println(reflect.TypeOf(stepsValue))
	fmt.Println(output["output"].String())
	outputDictAsMap := convert.FromStringDict(output)
	outputObject := outputDictAsMap["output"]
	outputTemplate, err := yaml.Marshal(&outputObject)

	return string(outputTemplate), err
}

func handlePanic() {
	if error := recover(); error != nil {
		log.Error("Recovering from panic: ", error)
	}
}

// Function map for 'vela' function
func VelaFuncMap() template.FuncMap {
	return template.FuncMap(map[string]interface{}{
		"vela": vela,
	})
}

// Simulates the 'vela' function during validation, so as to not fail templates that use it
func vela(variableName string) (envVariable string, err error) {
	if variableName != "" {
		envVariable = "${" + strings.ToUpper(variableName) + "}"
	} else {
		err = errors.New("Environment variable name cannot be empty in 'vela' function")
	}

	return
}

func recursiveToValue(input interface{}) (out starlark.Value, err error) {
	if err != nil {
		return nil, err
	}
	switch input.(type) {
	case map[string]interface{}:
		dict := starlark.Dict{}
		for k, v := range input.(map[string]interface{}) {
			key, err := convert.ToValue(k)
			if err != nil {
				return nil, err
			}
			val, err := recursiveToValue(v)
			if err != nil {
				return nil, err
			}
			dict.SetKey(key, val)
		}
		return &dict, nil
	case []interface{}:
		l := input.([]interface{})
		out := make([]starlark.Value, 0, len(l))
		for i := 0; i < len(l); i++ {
			val, err := recursiveToValue(out[i])
			if err != nil {
				return nil, err
			}
			out[i] = val
		}
		return starlark.NewList(out), nil
	default:
		return convert.ToValue(input)
	}
}

func MakeStarlarkStringDict(m map[string]interface{}) (starlark.StringDict, error) {
	dict := make(starlark.StringDict, len(m))
	for k, v := range m {
		val, err := recursiveToValue(v)
		if err != nil {
			return nil, err
		}
		dict[k] = val
	}
	return dict, nil
}