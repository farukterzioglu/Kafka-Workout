package commands

// RateReviewCommand command for rating a review
type RateReviewCommand struct {
	ReviewID int32 `json:"reviewId"`
	Star     int8  `json:"star"`
}
