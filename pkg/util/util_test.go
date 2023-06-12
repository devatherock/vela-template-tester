//go:build test
// +build test

package util

import (
	"testing"

	"github.com/devatherock/vela-template-tester/test/helper"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogLevel(test *testing.T) {
	helper.SetEnvironmentVariable(test, "PARAMETER_LOG_LEVEL", "debug")

	InitLogLevel()
	assert.Equal(test, "debug", log.GetLevel().String())
}

func TestInitWithInvalidLogLevel(test *testing.T) {
	helper.SetEnvironmentVariable(test, "PARAMETER_LOG_LEVEL", "hola")

	InitLogLevel()
	assert.Equal(test, "info", log.GetLevel().String())
}
