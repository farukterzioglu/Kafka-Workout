package commandhandlers

import "github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"

// HandlerRequest request model for handlers
type HandlerRequest struct {
	Command         commands.ICommand
	HandlerResponse chan<- interface{}
}

// ICommandHandler interface for command handlers
type ICommandHandler interface {
	HandleAsync(request HandlerRequest)
}
