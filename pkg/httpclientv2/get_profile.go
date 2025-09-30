package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_client"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

func (c *Client) Get(s string) (*http_client.GetProfileOutput, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("uuid.Parse: %w", domain.ErrUUIDInvalid)
	}

	output, err := c.client.GetProfileByIDWithResponse(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("GetProfileByIdWithResponse: %w", err)
	}

	if output.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if output.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request failed: status: %s, body:%s", output.Status(), output.Body)
	}

	return output.JSON200, nil
}
