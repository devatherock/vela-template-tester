package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/qri-io/starlib"
	log "github.com/sirupsen/logrus"
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

func Validate(validationRequest ValidationRequest) (validationResponse ValidationResponse) {
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
	parsedTemplate, _ := template.New("test").Funcs(VelaFuncMap()).Funcs(sprig.TxtFuncMap()).Parse(validationRequest.Template)
	err := parsedTemplate.Execute(buffer, validationRequest.Parameters)

	outputTemplate := buffer.String()
	return outputTemplate, err
}

func validateStarlarkTemplate(validationRequest *ValidationRequest) (string, error) {
	context := map[string]interface{}{
		"vars": validationRequest.Parameters,
	}
	contextJson, _ := jsoniter.Marshal(context)

	requestBody := validationRequest.Template + "\nctx = " + string(contextJson) + "\nprint(main(ctx))"
	file, _ := os.CreateTemp("", "exec_starlark")
	os.WriteFile(file.Name(), []byte(requestBody), 0755)

	output := ""
	thread := &starlark.Thread{
		Print: func(thread *starlark.Thread, msg string) {
			output = msg
		},
		Load: starlib.Loader,
	}

	outputTemplate := make([]byte, 1)
	_, err := starlark.ExecFile(thread, file.Name(), nil, nil)
	if err == nil {
		outputObject := make(map[string]interface{})
		json.Unmarshal([]byte(output), &outputObject)
		outputTemplate, err = yaml.Marshal(&outputObject)
	}

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
		err = errors.New("environment variable name cannot be empty in 'vela' function")
	}

	return
}
