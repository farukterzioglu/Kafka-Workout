package main

import (
	_ "github.com/bsm/sarama-cluster"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/CommandHandlers"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Models"
)

func main() {
	var createReviewHandler *commandhandlers.CreateReviewHandler
	createReviewHandler = commandhandlers.NewCreateReviewHandler()

	review := models.Review{
		Text: "sample review",
	}

	createReviewHandler.HandleAsync(commands.CreateReviewCommand{
		Review: review,
	})

}
