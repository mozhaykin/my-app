package httpclient

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"
)

func (c *Client) Delete(ctx context.Context, id string) error {
	const deleteProfile = "amozhaykin/my-app/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, deleteProfile, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, http.NoBody)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close() //nolint:contextcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		switch {
		case resp.StatusCode == http.StatusNotFound:
			return ErrNotFound

		default:
			return fmt.Errorf("request failed: status: %s body: %s", resp.Status, body)
		}
	}

	return nil
}
