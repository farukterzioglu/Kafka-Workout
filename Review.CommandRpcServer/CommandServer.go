package main

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc/metadata"

	"github.com/farukterzioglu/KafkaComparer/Review.CommandEngine/Models"
	pb "github.com/farukterzioglu/KafkaComparer/Review.CommandRpcServer/reviewservice"
)

// CommandServer for handling rpc commands
type CommandServer struct {
}

// NewCommandServer creates and return a CommandServer instance
func NewCommandServer() *CommandServer {
	s := &CommandServer{}
	return s
}

// SaveReview handles SaveReview rpc command
func (server *CommandServer) SaveReview(ctx context.Context, request *pb.NewReviewRequest) (*pb.ReviewId, error) {
	review := models.Review{
		Text: request.Review.Text,
		Star: int8(request.Review.Star),
	}

	fmt.Printf("Got a new review : %s (%d) \n", review.Text, review.Star)

	// TODO : save the review

	return &pb.ReviewId{ReviewId: "0"}, nil
}

// SaveReviews handles SaveReviews rpc command
func (server *CommandServer) SaveReviews(stream pb.ReviewService_SaveReviewsServer) error {
	md, _ := metadata.FromIncomingContext(stream.Context())
	_ = md["batchCount"] // batchCount

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		var review models.Review
		review = models.Review{
			Text: request.Review.Text,
			Star: int8(request.Review.Star),
		}
		fmt.Printf("Received review with text : %s\n", review.Text)

		var reviewID string
		// TODO : Process review
		reviewID = "0000"

		if err := stream.Send(&pb.ReviewId{ReviewId: reviewID}); err != nil {
			return err
		}
	}
}

// GetTopReviews returns top 'GetTopReviewsRequest.count' reviews
func (server *CommandServer) GetTopReviews(req *pb.GetTopReviewsRequest, stream pb.ReviewService_GetTopReviewsServer) error {
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
