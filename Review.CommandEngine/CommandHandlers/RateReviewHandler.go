package commandhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
)

// RateReviewHandler is the handler for CreateReview command
type RateReviewHandler struct{}

// HandleAsync handles string message
func (handler *RateReviewHandler) HandleAsync(ctx context.Context, request HandlerRequest) {
	var rateReviewCommand commands.RateReviewCommand
	json.Unmarshal(request.Command, &rateReviewCommand)

	fmt.Printf("Review (%d) rated with star : %d", rateReviewCommand.ReviewID, rateReviewCommand.Star)

	request.HandlerResponse <- true
}

// NewRateReviewHandler creates and returns new 'rate review' command handler
func NewRateReviewHandler() *RateReviewHandler {
	return &RateReviewHandler{}
}
