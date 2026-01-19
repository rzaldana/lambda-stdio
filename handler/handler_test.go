package handler_test

import (
	//"fmt"
	"bufio"
	"context"
	//"encoding/json"
	"io"
	"os"
	//"strings"
	"syscall"
	"testing"
	"strconv"
	"time"

	"github.com/rzaldana/lambda-stdio/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//func TestNewHandlerFailsIfHandlerEnvVarIsNotSetAndHandlerIsNotPassedInAsOptionFunc(t *testing.T) {
//	_, err := handler.NewHandler()
//	require.Error(t, err)
//}
//
//func TestNewHandlerFailsIfHandlerIsNotAValidCommand(t *testing.T) {
//	_, err := handler.NewHandler(handler.WithHandlerCommand("./notvalid"))
//	require.Error(t, err)
//}
//
//func TestHandleRequestRunsHandlerCommandAndReturnsStdout (t *testing.T) {
//	h, err := handler.NewHandler(handler.WithHandlerCommand("./test_scripts/TestHandleRequestRunsHandlerCommand.bash"))
//	require.NoError(t, err)
//
//	response, err := h.HandleRequest(context.TODO(), []byte{})
//	require.NoError(t, err)
//	assert.Equal(t, string(json.RawMessage("hello world!")), 	string(response))
//}
//
//func TestHandleRequestWritesCommandsStderrOutputToStderrWriter (t *testing.T) {
//	var mockStderr strings.Builder 
//
//	h, err := handler.NewHandler(
//		handler.WithHandlerCommand("./test_scripts/TestHandleRequestWritesCommandsStderrOutputToStderrWriter.bash"),
//		handler.WithStderrWriter(&mockStderr),
//	)
//	require.NoError(t, err)
//
//	response, err := h.HandleRequest(context.TODO(), []byte{})
//	require.NoError(t, err)
//	assert.Equal(t, "stdout", string(response))
//	assert.Equal(t, "error!", mockStderr.String())
//}


func TestHandleRequestWritesCommandsStderrOutputToStderrWriterWheneverANewLineIsWritten (t *testing.T) {
	stderrReader, stderrWriter := io.Pipe()


	h, err := handler.NewHandler(
		handler.WithHandlerCommand("./test_scripts/TestHandleRequestWritesCommandsStderrOutputToStderrWriterWheneverANewLineIsWritten.bash"),
		handler.WithStderrWriter(stderrWriter),
	)
	require.NoError(t, err)

	// run handler with handlerCommand in parallel thread
	var response []byte 
	var handlerErr error
	ch := make(chan int)
	go func() {
		response, handlerErr = h.HandleRequest(context.TODO(), []byte{})
		ch <- 0
	}()



	// Create scanner to read newline delimited tokens from stderrReader
	scanner := bufio.NewScanner(stderrReader)

	// Read first line, which sould contain PID of handlerCommand
	scanner.Scan()
	handlerCommandPID, err := strconv.Atoi(scanner.Text())
	require.NoError(t, err)

	// wait 1 second and then Send signal to handlerCommand process to terminate
	// Need to wait a bit because otherwise the process dies before it's
	// even set traps
	time.Sleep(1 * time.Second)
	p, err := os.FindProcess(handlerCommandPID)
	require.NoError(t, err)
	p.Signal(syscall.SIGTERM)

	// wait for handlerCommand to return
	<-ch

	// Make sure handleRequest returned no errors
	require.NoError(t, handlerErr)

	// Verify stdout output was correct
	assert.Equal(t, "stdout", string(response))
}
