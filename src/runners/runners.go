package runners

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
)

var (
	Executable = run_executable
)

func run_executable(path string, writer http.ResponseWriter, request *http.Request, errorCallback func(err error)) {
	myPipeReader, handlerPipeWriter := io.Pipe()
	defer myPipeReader.Close()
	defer handlerPipeWriter.Close()

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		errorCallback(err)
		return
	}

	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewBuffer(requestBody)
	cmd.Stdout = handlerPipeWriter
	cmd.Stderr = handlerPipeWriter

	go func() {
		defer handlerPipeWriter.Close()

		if err := cmd.Run(); err != nil {
			errorCallback(err)
		}
	}()

	if _, err := io.Copy(writer, myPipeReader); err != nil {
		errorCallback(err)
	}
}
