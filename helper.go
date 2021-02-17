// +build integration

package main

import (
	"io/ioutil"
	"os/exec"
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
