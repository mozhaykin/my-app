package httpclientv2

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) Delete(id string) error {
	output, err := c.client.DeleteProfileByIDWithResponse(context.Background(), id)
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
