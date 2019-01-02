package commandhandlers

import (
	"encoding/json"
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
)

// CreateReviewHandler is the handler for CreateReview command
type CreateReviewHandler struct{}

// HandleAsync handles string message
func (handler *CreateReviewHandler) HandleAsync(request HandlerRequest) {
	var createReviewCommand commands.CreateReviewCommand
	json.Unmarshal(request.Command, &createReviewCommand)

	fmt.Println(string(request.Command[:]))

	fmt.Printf("Review create with text : %s \n", createReviewCommand.Review.Text)
	request.HandlerResponse <- true
}

// NewCreateReviewHandler creates and returns new 'create review' command handler
func NewCreateReviewHandler() *CreateReviewHandler {
	return &CreateReviewHandler{}
}
