package grpcclientv1

import (
	"context"
	"fmt"

	pb "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/grpc/profile_v1"
)

func (c *Client) Delete(ctx context.Context, id string) error {
	input := &pb.DeleteProfileInput{
		Id: id,
	}

	_, err := c.client.DeleteProfile(ctx, input)
	if err != nil {
		return fmt.Errorf("c.client.DeleteProfile: %w", err)
	}

	return nil
}
