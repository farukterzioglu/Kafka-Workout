package commandhandlers

import (
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
)

type CreateReviewHandler struct{}

func (handler *CreateReviewHandler) HandleAsync(command commands.CreateReviewCommand) {
	fmt.Printf("Review create with text : %s", command.Review.Text)
}

func NewCreateReviewHandler() *CreateReviewHandler {
	return &CreateReviewHandler{}
}
