package commandhandlers

import (
	"fmt"
)

// DefaultHandler is default command handler. Prints command text
type DefaultHandler struct{}

// HandleAsync handles string message
func (handler *DefaultHandler) HandleAsync(command string) {
	fmt.Printf("Command text : %s", command)
}

// NewDefaultHandler creates and returns new Default Handler
func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}
