//go:build !integration
// +build !integration

package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogLevel(test *testing.T) {
	setEnvironmentVariable(test, "PARAMETER_LOG_LEVEL", "debug")

	initLogLevel()
	assert.Equal(test, "debug", log.GetLevel().String())
}

func TestInitWithInvalidLogLevel(test *testing.T) {
	setEnvironmentVariable(test, "PARAMETER_LOG_LEVEL", "hola")

	initLogLevel()
	assert.Equal(test, "info", log.GetLevel().String())
}
