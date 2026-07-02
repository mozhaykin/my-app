package grpcclientv1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	pb "github.com/mozhaykin/my-app/gen/grpc/profile_v1"
)

type CreateProfileRequest struct {
	Name  string
	Age   int
	Email string
	Phone string
}

func (c *Client) Create(ctx context.Context, r CreateProfileRequest) (uuid.UUID, error) {
	input := &pb.CreateProfileInput{
		Name:  r.Name,
		Age:   int32(r.Age), //nolint: gosec
		Email: r.Email,
		Phone: r.Phone,
	}

	resp, err := c.client.CreateProfile(ctx, input)
	if err != nil {
		return uuid.Nil, fmt.Errorf("c.client.CreateProfile: %w", err)
	}

	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("uuid.Parse: %w", err)
	}

	return id, nil
}
