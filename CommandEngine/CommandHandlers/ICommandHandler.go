package commandhandlers

// HandlerRequest request model for handlers
type HandlerRequest struct {
	Command         []byte
	HandlerResponse chan<- interface{}
}

// ICommandHandler interface for command handlers
type ICommandHandler interface {
	HandleAsync(request HandlerRequest)
}
