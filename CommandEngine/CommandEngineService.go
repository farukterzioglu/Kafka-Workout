package main

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/CommandHandlers"
)

// CommandRequest is the request type for commands
type CommandRequest struct {
	Msg        *sarama.ConsumerMessage
	ResponseCh chan interface{}
}

type commandCreatorFunc func() commandhandlers.ICommandHandler

var commandMap map[string]commandCreatorFunc

// CommandEngineService is service that handles command messages
type CommandEngineService struct{}

// NewCommandEngineService returns new command engine service
func NewCommandEngineService() *CommandEngineService {
	commandMap = make(map[string]commandCreatorFunc)
	commandMap["create-review"] = func() commandhandlers.ICommandHandler {
		return commandhandlers.NewCreateReviewHandler()
	}
	commandMap["rate-review"] = func() commandhandlers.ICommandHandler {
		return commandhandlers.NewRateReviewHandler()
	}

	return &CommandEngineService{}
}

// HandleMessage handles consumed command message
func (service *CommandEngineService) HandleMessage(request CommandRequest) {
	msg := request.Msg
	fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

	// Request
	var handlerRequest commandhandlers.HandlerRequest
	handlerRequest = commandhandlers.HandlerRequest{
		Command:         string(msg.Value[:]),
		HandlerResponse: make(chan interface{}),
	}

	// Handler
	var handler commandhandlers.ICommandHandler
	if createHandler, ok := commandMap[msg.Topic]; ok {
		handler = createHandler()
	} else {
		handler = commandhandlers.NewDefaultHandler()
	}
	handler.HandleAsync(handlerRequest)

	request.ResponseCh <- handlerRequest.HandlerResponse
}
