package commandhandlers

import (
	"context"
	"fmt"
)

// DefaultHandler is default command handler. Prints command text
type DefaultHandler struct{}

// HandleAsync handles string message
func (handler *DefaultHandler) HandleAsync(ctx context.Context, request HandlerRequest) {
	command := request.Command
	fmt.Printf("Command text : %s", command)

	request.HandlerResponse <- true
}

// NewDefaultHandler creates and returns new Default Handler
func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}
