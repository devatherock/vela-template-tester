package helper

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// Executes a command and returns its output and exit code
func ExecuteCommand(command *exec.Cmd) (int, string) {
	stdOut, _ := command.StdoutPipe()
	stdErr, _ := command.StderrPipe()
	command.Start()
	stdOutBytes, _ := ioutil.ReadAll(stdOut)
	stdErrBytes, _ := ioutil.ReadAll(stdErr)
	command.Wait()

	output := string(stdOutBytes) + string(stdErrBytes)
	return command.ProcessState.ExitCode(), output
}

// Sets an environment variable that will be cleaned up when the test ends
func SetEnvironmentVariable(test *testing.T, variable string, value string) {
	os.Setenv(variable, value)

	test.Cleanup(func() {
		os.Unsetenv(variable)
	})
}

func GetProjectRoot() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath) // Returns the absolute path to helper.go

	return filepath.Join(currentDir, "../..") // Path 2 levels above helper.go
}

// Gets absolute file path from relative path
func AbsolutePath(relativePath string) string {
	return filepath.Join(GetProjectRoot(), relativePath)
}

// Returns a relative or absolute path based on the flag
func Path(relativePath string, convertToAbsolute bool) string {
	if convertToAbsolute {
		return AbsolutePath(relativePath)
	} else {
		return relativePath
	}
}

// Returns the environment variable's value if it is present, else the default
func Getenv(variableName string, defaultValue string) string {
	value := os.Getenv(variableName)
	output := defaultValue

	if value != "" {
		output = value
	}

	return output
}

// Returns true if the input string equals 'true', returns false otherwise
func StringToBool(value string) bool {
	boolValue, err := strconv.ParseBool(value)

	if err != nil {
		boolValue = false
	}

	return boolValue
}

// Converts a string command into a command string and an arguments array
func ParseCommand(inputCommand string) (string, []string) {
	parts := strings.Split(inputCommand, " ")
	command := parts[0]

	arguments := []string{}
	for index, part := range parts {
		if index > 0 {
			arguments = append(arguments, part)
		}
	}

	return command, arguments
}
