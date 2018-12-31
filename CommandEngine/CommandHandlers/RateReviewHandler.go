package commandhandlers

import (
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Commands"
)

// RateReviewHandler is the handler for CreateReview command
type RateReviewHandler struct{}

// HandleAsync handles string message
func (handler *RateReviewHandler) HandleAsync(request HandlerRequest) {
	command := request.Command
	rateReviewCommand := command.(commands.RateReviewCommand)

	fmt.Printf("Review (%d) rated with star : %d", rateReviewCommand.ReviewID, rateReviewCommand.Star)

	request.HandlerResponse <- true
}

// NewRateReviewHandler creates and returns new 'rate review' command handler
func NewRateReviewHandler() *RateReviewHandler {
	return &RateReviewHandler{}
}
