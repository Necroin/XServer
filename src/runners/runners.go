package runners

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
)

func Executable(path string, writer io.Writer, request io.Reader, errorCallback func(string, error), logCallback func(string), args ...string) {
	myPipeReader, handlerPipeWriter := io.Pipe()
	defer myPipeReader.Close()
	defer handlerPipeWriter.Close()

	requestBody, err := ioutil.ReadAll(request)
	if err != nil {
		errorCallback("failed read request body", err)
		return
	}

	cmd := exec.Command(path, args...)
	cmd.Stdin = bytes.NewBuffer(requestBody)
	cmd.Stdout = handlerPipeWriter
	cmd.Stderr = handlerPipeWriter

	go func() {
		defer handlerPipeWriter.Close()
		logCallback("run file")
		if err := cmd.Run(); err != nil {
			errorCallback("failed run handler file", err)
		}
	}()

	if _, err := io.Copy(writer, myPipeReader); err != nil {
		errorCallback("failed copy handler response", err)
	}
}

func Tool(tool string, path string, writer io.Writer, request io.Reader, errorCallback func(string, error), logCallback func(string), args ...string) {
	cmdArgs := append([]string{path}, args...)
	Executable(tool, writer, request, errorCallback, logCallback, cmdArgs...)
}

func Python(path string, writer io.Writer, request io.Reader, errorCallback func(string, error), logCallback func(string), args ...string) {
	Tool("python", path, writer, request, errorCallback, logCallback, args...)
}

func Lua(path string, writer io.Writer, request io.Reader, errorCallback func(string, error), logCallback func(string), args ...string) {
	Tool("lua", path, writer, request, errorCallback, logCallback, args...)
}
