package grpcclientv1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	pb "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/grpc/profile_v1"
)

type Profile struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Status    int       `json:"status"`
	Verified  bool      `json:"verified"`
	Contacts  struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contacts"`
}

func (c *Client) Get(id string) (Profile, error) {
	input := &pb.GetProfileInput{
		Id: id,
	}

	output, err := c.client.GetProfile(context.Background(), input)
	if err != nil {
		return Profile{}, fmt.Errorf("c.client.GetProfile: %w", err)
	}

	return Profile{
		ID:        uuid.MustParse(output.GetId()),
		CreatedAt: output.GetCreatedAt().AsTime(),
		UpdatedAt: output.GetUpdatedAt().AsTime(),
		Name:      output.GetName(),
		Age:       int(output.GetAge()),
		Status:    int(output.GetStatus()),
		Verified:  output.GetVerified(),
		Contacts: struct {
			Email string `json:"email"`
			Phone string `json:"phone"`
		}{
			Email: output.GetContacts().GetEmail(),
			Phone: output.GetContacts().GetPhone(),
		},
	}, nil
}
