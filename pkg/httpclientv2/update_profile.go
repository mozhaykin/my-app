package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/client"
)

type UpdateProfileRequest struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

func (c *Client) Update(ctx context.Context, request UpdateProfileRequest) error {
	input := client.UpdateProfileInput{
		ID:    request.ID,
		Name:  request.Name,
		Age:   request.Age,
		Email: request.Email,
		Phone: request.Phone,
	}

	output, err := c.client.UpdateProfileWithResponse(ctx, input)
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
