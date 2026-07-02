package httpclientv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mozhaykin/my-app/gen/http/profile_v2/client"
)

type GetProfilesRequest struct {
	Sort   string
	Order  string
	Offset int
	Limit  int
}

func (c *Client) GetProfiles(ctx context.Context, r GetProfilesRequest) ([]client.GetProfileOutput, error) {
	params := client.GetProfilesParams{
		Sort:   r.Sort,
		Order:  &r.Order,
		Offset: &r.Offset,
		Limit:  &r.Limit,
	}

	output, err := c.client.GetProfilesWithResponse(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("GetProfilesWithResponse: %w", err)
	}

	if output.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if output.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request failed: status: %s, body:%s", output.Status(), output.Body)
	}

	return *output.JSON200, nil
}
