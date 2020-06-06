package main

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "os"
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

func init() {
    level, ok := os.LookupEnv("PARAMETER_LOG_LEVEL")

    // LOG_LEVEL not set, default to info
    if !ok {
        level = "info"
    }

    // parse string, this is built-in feature of logrus
    logLevel, err := log.ParseLevel(level)
    if err != nil {
        logLevel = log.InfoLevel
    }

    // set global log level
    log.SetLevel(logLevel)
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