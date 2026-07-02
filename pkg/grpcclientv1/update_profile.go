package grpcclientv1

import (
	"context"
	"fmt"

	pb "github.com/mozhaykin/my-app/gen/grpc/profile_v1"
)

type UpdateProfileRequest struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

func (c *Client) Update(ctx context.Context, r UpdateProfileRequest) error {
	input := &pb.UpdateProfileInput{
		Id:    r.ID,
		Name:  r.Name,
		Age:   parseAge(r.Age),
		Email: r.Email,
		Phone: r.Phone,
	}

	_, err := c.client.UpdateProfile(ctx, input)
	if err != nil {
		return fmt.Errorf("c.client.UpdateProfile: %w", err)
	}

	return nil
}

func parseAge(age *int) *int32 {
	if age == nil {
		return nil
	}

	a := int32(*age) //nolint: gosec

	return &a
}
