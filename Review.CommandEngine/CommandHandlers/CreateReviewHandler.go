package commandhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Models"
	"github.com/google/uuid"
)

// CreateReviewHandler is the handler for CreateReview command
type CreateReviewHandler struct{}

// HandleAsync handles string message
func (handler *CreateReviewHandler) HandleAsync(ctx context.Context, request HandlerRequest) {
	var createReviewCommand commands.CreateReviewCommand
	json.Unmarshal(request.Command, &createReviewCommand)

	fmt.Println(string(request.Command[:]))

	fmt.Printf("Review create with text : %s \n", createReviewCommand.Review.Text)

	// TODO : Store review

	// TODO : Get unique id. insterad for now create uuid
	var reviewID string
	uuid, err := uuid.NewRandom()
	if err != nil {
		reviewID = ""
	} else {
		reviewID = string(uuid[:])
	}

	ctx = models.NewContextWithReviewId(ctx, reviewID)

	request.HandlerResponse <- true
}

// NewCreateReviewHandler creates and returns new 'create review' command handler
func NewCreateReviewHandler() *CreateReviewHandler {
	return &CreateReviewHandler{}
}
