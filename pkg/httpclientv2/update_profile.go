package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/client"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

type UpdateProfileRequest struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

func (c *Client) Update(request UpdateProfileRequest) error {
	id, err := uuid.Parse(request.ID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", domain.ErrUUIDInvalid)
	}

	input := client.UpdateProfileInput{
		ID:    id,
		Name:  request.Name,
		Age:   request.Age,
		Email: request.Email,
		Phone: request.Phone,
	}

	output, err := c.client.UpdateProfileWithResponse(context.Background(), input)
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	if output.StatusCode() == http.StatusNotFound {
		return ErrNotFound
	}

	if output.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("request failed: status: %s, body:%s", output.Status(), output.Body)
	}

	return nil
}
