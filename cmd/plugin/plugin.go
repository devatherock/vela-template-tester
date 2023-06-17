package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/devatherock/vela-template-tester/pkg/util"
	"github.com/devatherock/vela-template-tester/pkg/validator"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type PluginValidationRequest struct {
	InputFile      string                 `json:"input_file,omitempty"`
	Variables      map[string]interface{} `json:",omitempty"`
	ExpectedOutput string                 `json:"expected_output,omitempty"`
	TemplateType   string                 `json:"template_type,omitempty"`
}

var exit func(code int) = os.Exit

// Initializes log level
func init() {
	util.InitLogLevel()
}

// Plugin entry point
func main() {
	runApp(os.Args)
}

// Reads the plugin parameters and runs it
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
			Name:    "template-type",
			Aliases: []string{"tt"},
			Usage:   "The template type. Needs to be 'starlark' if '--input-file' is a starlark template",
			EnvVars: []string{"TEMPLATE_TYPE", "PARAMETER_TEMPLATE_TYPE"},
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
	util.HandleError(err)
}

// Tests the supplied templates using the validator
func run(context *cli.Context) error {
	pluginValidationRequests := readInputParameters(context)
	var validationFailure bool
	var validationStatus error // For easier testing

	for _, request := range pluginValidationRequests {
		validationRequest := validator.ValidationRequest{}

		content, error := os.ReadFile(request.InputFile)
		if error != nil {
			return error
		}
		validationRequest.Template = string(content)
		validationRequest.Parameters = request.Variables
		validationRequest.Type = request.TemplateType

		validationResponse := validator.Validate(validationRequest)
		if validationResponse.Error != "" {
			message := fmt.Sprintf("Template '%s' is invalid. Error: %s", request.InputFile, validationResponse.Error)
			validationStatus = errors.New(message)

			log.Error(message)
			validationFailure = true
		} else {
			validationResult := verifyOutput(request, validationResponse)

			if validationResult {
				log.Printf("Template '%s' is valid.", request.InputFile)
			} else {
				message := fmt.Sprintf("Template '%s' is valid, but did not match expected output", request.InputFile)
				validationStatus = errors.New(message)

				log.Error(message)
				validationFailure = true
			}
		}
	}

	if validationFailure {
		exit(1)
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
			util.HandleError(error)

			pluginValidationRequest.Variables = parsedVariables
		}

		expectedOutputFile := context.String("expected-output")
		if expectedOutputFile != "" {
			pluginValidationRequest.ExpectedOutput = expectedOutputFile
		}

		templateType := context.String("template-type")
		if templateType != "" {
			pluginValidationRequest.TemplateType = templateType
		}

		pluginValidationRequests = append(pluginValidationRequests, pluginValidationRequest)
	}

	// Parse array of validation requests if specified
	templates := context.String("templates")
	if templates != "" {
		suppliedValidationRequests := []PluginValidationRequest{}
		error := json.Unmarshal([]byte(templates), &suppliedValidationRequests)
		util.HandleError(error)

		pluginValidationRequests = append(pluginValidationRequests, suppliedValidationRequests...)
	}

	if len(pluginValidationRequests) == 0 {
		log.Warn("No template specified")
		exit(0)
	}

	return pluginValidationRequests
}

// Verifies if the processed template matches the expected output
func verifyOutput(request PluginValidationRequest, validationResponse validator.ValidationResponse) bool {
	if request.ExpectedOutput != "" {
		expectedOutput, error := os.ReadFile(request.ExpectedOutput)
		util.HandleError(error)

		expectedOutputMap := make(map[interface{}]interface{})
		error = yaml.Unmarshal([]byte(expectedOutput), &expectedOutputMap)
		util.HandleError(error)

		processedTemplateMap := make(map[interface{}]interface{})
		error = yaml.Unmarshal([]byte(validationResponse.Template), &processedTemplateMap)
		util.HandleError(error)

		return reflect.DeepEqual(expectedOutputMap, processedTemplateMap)
	}

	return true
}
