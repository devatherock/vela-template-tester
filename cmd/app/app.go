//go:build !plugin && !integration
// +build !plugin,!integration

package main

import (
	"io"
	"net/http"
	"os"

	"github.com/devatherock/vela-template-tester/pkg/util"
	"github.com/devatherock/vela-template-tester/pkg/validator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Initializes log level
func init() {
	util.InitLogLevel()
}

func main() {
	http.HandleFunc("/api/expandTemplate", expandTemplate)
	http.HandleFunc("/api/health", checkHealth)

	http.ListenAndServe(":"+lookupPort(), nil)
}

// Handles /api/expandTemplate endpoint. Expands supplied template with
// the supplied parameters, to verify if it is valid
func expandTemplate(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/x-yaml")

	// Read request
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return
	}

	// Parse request
	validationRequest := validator.ValidationRequest{}
	yaml.Unmarshal(requestBody, &validationRequest)

	// Validate template
	validationResponse := validator.Validate(validationRequest)

	// Write response
	responseBody, err := yaml.Marshal(&validationResponse)
	if err != nil {
		log.Error("error: ", err)
	} else {
		writer.Write(responseBody)
	}
}

// Handles /api/health endpoint. Indicates the health of the application
func checkHealth(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("UP"))
}

// Reads port from PORT environment variable
func lookupPort() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	return port
}
