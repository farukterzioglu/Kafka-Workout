package main

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
)

type CommandEngineService struct{}

func NewCommandEngineService() *CommandEngineService {
	return &CommandEngineService{}
}

func (service *CommandEngineService) HandleMessage(msg *sarama.ConsumerMessage) {
	fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

	// var createReviewHandler *commandhandlers.CreateReviewHandler // TODO : ??? *commandhandlers.ICommandHandler
	// createReviewHandler = commandhandlers.NewCreateReviewHandler()

	// review := models.Review{
	// 	Text: "sample review",
	// }

	// createReviewHandler.HandleAsync(commands.CreateReviewCommand{
	// 	Review: review,
	// })
}
