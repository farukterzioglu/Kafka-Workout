package main

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/CommandHandlers"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Models"
)

// CommandRequest is the request type for commands
type CommandRequest struct {
	Msg        *sarama.ConsumerMessage
	ResponseCh chan interface{}
}

// CommandEngineService is service that handles command messages
type CommandEngineService struct{}

// NewCommandEngineService returns new command engine service
func NewCommandEngineService() *CommandEngineService {
	return &CommandEngineService{}
}

// HandleMessage handles consumed command message
func (service *CommandEngineService) HandleMessage(request CommandRequest) {
	msg := request.Msg
	fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

	var handler commandhandlers.ICommandHandler
	var cmd commands.ICommand

	type commandCreatorFunc func() commandhandlers.ICommandHandler
	commandMap := make(map[string]commandCreatorFunc)
	commandMap["create-review"] = func() commandhandlers.ICommandHandler {
		cmd = commands.CreateReviewCommand{
			Review : models.Review{
				Text : msg.Value
			}
		}
		return commandhandlers.NewCreateReviewHandler()
	}
	commandMap["rate-review"] = func() commandhandlers.ICommandHandler {
		cmd = commands.RateReviewCommand{}
		return commandhandlers.NewRateReviewHandler()
	}

	ok, handler = commandMap[msgh.Topic]()
	if !ok {
		handler = commandhandlers.NewDefaultHandler()
	}

	var handlerRequest commandhandlers.HandlerRequest
	handlerRequest = commandhandlers.HandlerRequest{
		Command:         cmd,
		HandlerResponse: make(chan interface{}),
	}
	handler.HandleAsync(handlerRequest)

	request.ResponseCh <- handlerRequest.HandlerResponse
}
