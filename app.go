// +build !plugin,!integration

package main

import (
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	// Parse request
	validationRequest := ValidationRequest{}
	err = yaml.Unmarshal(requestBody, &validationRequest)

	// Validate template
	validationResponse := validate(validationRequest)

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

// Reads port from PORT environment variable available on heroku
func lookupPort() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	return port
}
