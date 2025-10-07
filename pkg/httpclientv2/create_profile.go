package httpclientv2

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/client"
)

type CreateProfileRequest struct {
	Name  string
	Age   int
	Email string
	Phone string
}

func (c *Client) Create(r CreateProfileRequest) (uuid.UUID, error) {
	input := client.CreateProfileInput{
		Name:  r.Name,
		Age:   r.Age,
		Email: openapi_types.Email(r.Email),
		Phone: r.Phone,
	}

	output, err := c.client.CreateProfileWithResponse(context.Background(), input)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create profile: %w", err)
	}

	if output.StatusCode() != http.StatusCreated {
		return uuid.Nil, fmt.Errorf("create profile: %w", errors.New(output.JSON400.Error))
	}

	return output.JSON201.ID, nil
}
