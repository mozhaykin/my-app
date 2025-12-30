package httpclientv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type GetProfilesRequest struct {
	Sort   string
	Order  string
	Offset int
	Limit  int
}

type GetProfilesResponse struct {
	Profiles []Profile `json:"profiles"`
}

func (c *Client) GetProfiles(ctx context.Context, request GetProfilesRequest) ([]Profile, error) {
	u := &url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   "amozhaykin/my-app/api/v1/profiles",
	}

	// Query params
	q := u.Query()
	q.Set("sort", request.Sort)
	q.Set("order", request.Order)
	q.Set("offset", strconv.Itoa(request.Offset))
	q.Set("limit", strconv.Itoa(request.Limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: status: %s body: %s", resp.Status, body)
	}

	var out GetProfilesResponse

	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return out.Profiles, nil
}
