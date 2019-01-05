package commandhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/farukterzioglu/KafkaComparer/Review.CommandEngine/Commands"
	"github.com/farukterzioglu/KafkaComparer/Review.CommandEngine/Models"
	pb "github.com/farukterzioglu/KafkaComparer/Review.CommandRpcServer/reviewservice"
	"google.golang.org/grpc/metadata"
)

// CreateReviewHandler is the handler for CreateReview command
type CreateReviewHandler struct {
	client pb.ReviewServiceClient
}

// NewCreateReviewHandler creates and returns new 'create review' command handler
func NewCreateReviewHandler(c pb.ReviewServiceClient) *CreateReviewHandler {
	return &CreateReviewHandler{
		client: c,
	}
}

func populateRPCCommand(review *models.Review) *pb.NewReviewRequest {
	return &pb.NewReviewRequest{
		Review: &pb.Review{
			Text: review.Text,
			Star: int32(review.Star),
		},
	}
}

// HandleAsync handles string message
func (handler *CreateReviewHandler) HandleAsync(ctx context.Context, request HandlerRequest) {
	fmt.Println(string(request.Command[:]))

	var createReviewCommand commands.CreateReviewCommand
	json.Unmarshal(request.Command, &createReviewCommand)

	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("request-owner", "CreateReviewHandler"),
	)
	// metadata.AppendToOutgoingContext(ctx, "key", "value")

	response, err := handler.client.SaveReview(ctx, populateRPCCommand(&createReviewCommand.Review))
	if err != nil {
		request.ErrResponse <- err
		return
	}

	fmt.Printf("Review create with id : %s \n", response.ReviewId)

	ctx = models.NewContextWithReviewId(ctx, response.ReviewId)
	request.HandlerResponse <- response.ReviewId
}
