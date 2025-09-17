package httpclientv2

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_client"
)

type CreateProfileRequest struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (c *Client) Create(request CreateProfileRequest) (uuid.UUID, error) {
	input := http_client.CreateProfileInput{
		Name:  request.Name,
		Age:   request.Age,
		Email: openapi_types.Email(request.Email),
		Phone: request.Phone,
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
