package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

// Executes a command and returns its output and exit code
func executeCommand(command *exec.Cmd) (int, string) {
	stdOut, _ := command.StdoutPipe()
	stdErr, _ := command.StderrPipe()
	command.Start()
	stdOutBytes, _ := ioutil.ReadAll(stdOut)
	stdErrBytes, _ := ioutil.ReadAll(stdErr)
	command.Wait()

	output := string(stdOutBytes) + string(stdErrBytes)
	return command.ProcessState.ExitCode(), output
}

func setEnvironmentVariable(test *testing.T, variable string, value string) {
	os.Setenv(variable, value)

	test.Cleanup(func() {
		os.Unsetenv(variable)
	})
}
