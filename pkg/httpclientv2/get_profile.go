package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_client"
)

func (c *Client) Get(id string) (*http_client.GetProfileOutput, error) {
	output, err := c.client.GetProfileByIDWithResponse(context.Background(), uuid.MustParse(id))
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
