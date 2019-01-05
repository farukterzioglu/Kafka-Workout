package commands

import "github.com/farukterzioglu/KafkaComparer/Review.CommandEngine/Models"

type CreateReviewCommand struct {
	Review models.Review `json:"review"`
}
