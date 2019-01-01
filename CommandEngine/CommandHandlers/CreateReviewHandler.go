package commandhandlers

import (
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Models"
)

// CreateReviewHandler is the handler for CreateReview command
type CreateReviewHandler struct{}

// HandleAsync handles string message
func (handler *CreateReviewHandler) HandleAsync(request HandlerRequest) {
	// TODO : Parse json
	createReviewCommand := commands.CreateReviewCommand{
		Review: models.Review{
			Text: string(request.Command[:]),
		},
	}

	fmt.Printf("Review create with text : %s \n", createReviewCommand.Review.Text)

	request.HandlerResponse <- true
}

// NewCreateReviewHandler creates and returns new 'create review' command handler
func NewCreateReviewHandler() *CreateReviewHandler {
	return &CreateReviewHandler{}
}
