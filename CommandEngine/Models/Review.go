package models

import "context"

type key string

const reviewIDKey key = "reviewIDKey"

// Review struct
type Review struct {
	Text string `json:"text"`
	Star int    `json:"star"`
}

// NewContext returns a new Context that carries a provided review id value
func NewContextWithReviewId(ctx context.Context, reviewID string) context.Context {
	return context.WithValue(ctx, reviewIDKey, reviewID)
}

// FromContext extracts a review id from a Context
func ReviewIdFromContext(ctx context.Context) string {
	return ctx.Value(reviewIDKey).(string)
}
