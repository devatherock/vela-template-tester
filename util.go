//go:build !integration
// +build !integration

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Initializes log level
func init() {
	initLogLevel()
}

func initLogLevel() {
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

// Logs the error and exits the application
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
