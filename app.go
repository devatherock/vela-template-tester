package main

import (
    "io/ioutil"
    "net/http"
    "os"

    "gopkg.in/yaml.v2"
    log "github.com/sirupsen/logrus"
)

// Initializes log level
func init() {
    initLogLevel()
}

func main() {
    http.HandleFunc("/api/expandTemplate", func(writer http.ResponseWriter, request *http.Request) {
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
    })

    http.HandleFunc("/api/health", func(writer http.ResponseWriter, request *http.Request) {
        writer.Write([]byte("UP"))
    })

    // Read from PORT environment variable available on heroku
    port, ok := os.LookupEnv("PORT")
    if !ok {
        port = "8080"
    }
    http.ListenAndServe(":" + port, nil)
}
