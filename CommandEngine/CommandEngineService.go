package main

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/CommandHandlers"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
)

// CommandEngineService is service that handles command messages
type CommandEngineService struct{}

// NewCommandEngineService returns new command engine service
func NewCommandEngineService() *CommandEngineService {
	return &CommandEngineService{}
}

// HandleMessage handles consumed command message
func (service *CommandEngineService) HandleMessage(msg *sarama.ConsumerMessage) {
	fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

	var cmd commands.ICommand
	// TODO : deserialize json message
	cmd = commands.CreateReviewCommand{}

	switch parsedCmd := cmd.(type) {
	case commands.CreateReviewCommand:
		handler := commandhandlers.NewCreateReviewHandler()
		handler.HandleAsync(parsedCmd)
	default:
		// t is some other type that we didn't name.
	}
}
