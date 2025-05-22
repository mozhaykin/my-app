package httpclient

import (
	"fmt"
	"net/http"
)

func (c *Client) Delete(id string) error {
	const deleteProfile = "amozhaykin/my-app/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, deleteProfile, id)

	req, err := http.NewRequest(http.MethodDelete, path, http.NoBody)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request failed: status: %s", resp.Status)
	}

	return nil
}
