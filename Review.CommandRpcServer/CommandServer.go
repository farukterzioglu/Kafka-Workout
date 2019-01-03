package commandrpcserver

import (
	"context"
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/Review.CommandEngine/Models"
	pb "github.com/farukterzioglu/KafkaComparer/Review.CommandRpcServer/reviewservice"
)

// CommandServer for handling rpc commands
type CommandServer struct {
}

// SaveReview handles SaveReview rpc command
func (server *CommandServer) SaveReview(ctx context.Context, review *pb.Review) (*pb.ReviewId, error) {
	fmt.Printf("Got a new review : %s (%d) \n", review.Text, review.Star)

	// TODO : save the review

	return &pb.ReviewId{ReviewId: 0}, nil
}

// GetTopReviews returns top 'GetTopReviewsRequest.count' reviews
func (server *CommandServer) GetTopReviews(ctx context.Context, req *pb.GetTopReviewsRequest, stream pb.ReviewService_GetTopReviewsServer) error {
	// TODO : Get reviews
	reviewList := []models.Review{
		models.Review{
			Text: "First awesome review",
			Star: 5,
		},
		models.Review{
			Text: "Last meh review",
			Star: 3,
		},
	}

	for _, review := range reviewList {
		reviewReq := pb.Review{
			Text: review.Text,
			Star: int32(review.Star),
		}

		if err := stream.Send(&reviewReq); err != nil {
			return err
		}
	}

	return nil
}
