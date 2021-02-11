// +build !plugin

package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TODO merge test with TestMain once it is added
func TestInit(test *testing.T) {
	setEnvironmentVariable(test, "PARAMETER_LOG_LEVEL", "debug")

	initLogLevel()
	assert.Equal(test, "debug", log.GetLevel().String())
}
