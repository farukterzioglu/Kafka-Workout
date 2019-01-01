package commandhandlers

import (
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
)

// CreateReviewHandler is the handler for CreateReview command
type CreateReviewHandler struct{}

// HandleAsync handles string message
func (handler *CreateReviewHandler) HandleAsync(request HandlerRequest) {
	command := request.Command
	createReviewCommand := command.(commands.CreateReviewCommand)

	fmt.Printf("Review create with text : %s \n", createReviewCommand.Review.Text)

	request.HandlerResponse <- true
}

// NewCreateReviewHandler creates and returns new 'create review' command handler
func NewCreateReviewHandler() *CreateReviewHandler {
	return &CreateReviewHandler{}
}
