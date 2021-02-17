// +build !api,!integration

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type PluginValidationRequest struct {
	InputFile      string                 `json:"input_file,omitempty"`
	Variables      map[string]interface{} `json:",omitempty"`
	ExpectedOutput string                 `json:"expected_output,omitempty"`
}

// Plugin entry point
func main() {
	runApp(os.Args)
}

// Initializes and runs the app
func runApp(args []string) {
	app := cli.NewApp()
	app.Name = "vela template tester plugin"
	app.Action = run
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "input-file",
			Aliases: []string{"tf"},
			Usage:   "The template file to test",
			EnvVars: []string{"INPUT_FILE", "PARAMETER_INPUT_FILE"},
		},
		&cli.StringFlag{
			Name:    "templates",
			Aliases: []string{"ts"},
			Usage:   "The list of template files and variables to test",
			EnvVars: []string{"TEMPLATES", "PARAMETER_TEMPLATES"},
		},
		&cli.StringFlag{
			Name:    "variables",
			Aliases: []string{"v"},
			Usage:   "Variables to apply to the template",
			EnvVars: []string{"VARIABLES", "PARAMETER_VARIABLES"},
		},
		&cli.StringFlag{
			Name:    "expected-output",
			Aliases: []string{"o"},
			Usage:   "The expected output of the processed template",
			EnvVars: []string{"EXPECTED_OUTPUT", "PARAMETER_EXPECTED_OUTPUT"},
		},
	}

	err := app.Run(args)
	handleError(err)
}

// Tests the supplied templates using the validator
func run(context *cli.Context) error {
	pluginValidationRequests := readInputParameters(context)
	var validationStatus error

	for _, request := range pluginValidationRequests {
		validationRequest := ValidationRequest{}

		content, error := ioutil.ReadFile(request.InputFile)
		if error != nil {
			return error
		}
		validationRequest.Template = string(content)
		validationRequest.Parameters = request.Variables

		validationResponse := validate(validationRequest)
		if validationResponse.Error != "" {
			validationStatus = errors.New(fmt.Sprintf("Template '%s' is not valid. Error: %s", request.InputFile, validationResponse.Error))
		} else {
			validationResult := verifyOutput(request, validationResponse)

			if validationResult {
				log.Printf("Template '%s' is valid.", request.InputFile)
			} else {
				validationStatus = errors.New(fmt.Sprintf("Template '%s' is valid, but did not match expected output", request.InputFile))
			}
		}
	}

	return validationStatus
}

// Reads plugin input parameters
func readInputParameters(context *cli.Context) []PluginValidationRequest {
	pluginValidationRequests := []PluginValidationRequest{}

	// Create a plugin validation request out of the individual parameters
	templateFile := context.String("input-file")
	if templateFile != "" {
		pluginValidationRequest := PluginValidationRequest{
			InputFile: templateFile,
		}

		variables := context.String("variables")
		if variables != "" {
			parsedVariables := make(map[string]interface{})
			error := json.Unmarshal([]byte(variables), &parsedVariables)
			handleError(error)

			pluginValidationRequest.Variables = parsedVariables
		}

		expectedOutputFile := context.String("expected-output")
		if expectedOutputFile != "" {
			pluginValidationRequest.ExpectedOutput = expectedOutputFile
		}
		pluginValidationRequests = append(pluginValidationRequests, pluginValidationRequest)
	}

	// Parse array of validation requests if specified
	templates := context.String("templates")
	if templates != "" {
		suppliedValidationRequests := []PluginValidationRequest{}
		error := json.Unmarshal([]byte(templates), &suppliedValidationRequests)
		handleError(error)

		pluginValidationRequests = append(pluginValidationRequests, suppliedValidationRequests...)
	}

	if len(pluginValidationRequests) == 0 {
		log.Warn("No templates specified")
	}

	return pluginValidationRequests
}

// Verifies if the processed template matches the expected output
func verifyOutput(request PluginValidationRequest, validationResponse ValidationResponse) bool {
	if request.ExpectedOutput != "" {
		expectedOutput, error := ioutil.ReadFile(request.ExpectedOutput)
		handleError(error)

		expectedOutputMap := make(map[interface{}]interface{})
		error = yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)
		handleError(error)

		processedTemplateMap := make(map[interface{}]interface{})
		error = yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)
		handleError(error)

		return reflect.DeepEqual(expectedOutputMap, processedTemplateMap)
	}

	return true
}
