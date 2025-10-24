package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/client"
)

func (c *Client) Get(ctx context.Context, id string) (*client.GetProfileOutput, error) {
	output, err := c.client.GetProfileByIDWithResponse(ctx, id)
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
