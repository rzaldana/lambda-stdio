package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"bufio"
)


type config struct {
	handlerCommand *string	
	stderrWriter io.Writer
}

type Handler struct {
	config
}

type OptionFunction func(c *config)

func WithHandlerCommand(handlerCommand string) OptionFunction {
	return func(c *config) {
		c.handlerCommand = &handlerCommand	
	}
}

func WithStderrWriter(w io.Writer) OptionFunction {
	return func(c *config) {
		c.stderrWriter = w
	}
}

func NewHandler (optFns ...OptionFunction) (*Handler, error) {
	config := config{}
	for _, optFn := range optFns {
		optFn(&config)
	}

	if config.handlerCommand == nil {
		handlerEnvVarValue := os.Getenv("_HANDLER")
		if handlerEnvVarValue == "" {
			return nil, fmt.Errorf("_HANDLER environment variable not set and handlerCommand was not passed in config")
		}

		config.handlerCommand = &handlerEnvVarValue
	}

	// check if handlerCommand is valid
	_, err := exec.LookPath(*config.handlerCommand)
	if err != nil {
		return nil, fmt.Errorf("Command '%v' not available in PATH or in filesystem. Error: '%v'", *config.handlerCommand, err)
	}
	return &Handler{config: config}, nil
}

func (h *Handler) HandleRequest(ctx context.Context, event json.RawMessage) (json.RawMessage, error) {	
	cmd := exec.CommandContext(ctx, *h.handlerCommand)

	// Create stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to create stdout pipe. Error: '%v'", err)
	}

	// Create stderr pipe
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to create stdout pipe. Error: '%v'", err)
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Unable to start command '%v'. Error: '%v'", *h.handlerCommand, err)
	}

	// Read from stderr pipe and write to stderrWriter
	stderrChan := make(chan int)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			_, _ = h.stderrWriter.Write(append(scanner.Bytes(), '\n'))
			// TODO: handle this error
		}

		//if scanner.Err() != nil {
		//	
		//}
		// TODO: handler this error
		stderrChan <- 0
	}()

	// Read from stdout pipe
	var response json.RawMessage
	stdoutChan := make(chan int)
	go func() {
		response, err = io.ReadAll(stdout)
		stdoutChan <- 0
		//if err != nil {
		//	return nil, fmt.Errorf("Unable to read from stdout pipe. Error: '%v'", err) 
		//}
	}()

	// Read from stdout pipe
	// Will block until command handler's stdout is closed

	// Wait for command to exit
	err = cmd.Wait()
	// TODO: handle this error
	//if err != nil {
	//	fmt.Printf("Command failed: %v", err)
	//}
	<- stdoutChan
	return response, nil
}
