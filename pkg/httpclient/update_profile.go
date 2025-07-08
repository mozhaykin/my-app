package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Update(ctx context.Context, id, name string, age int, email, phone string) error {
	const updateProfile = "amozhaykin/my-app/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s", c.host, updateProfile)

	request := struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}{
		ID:    id,
		Name:  name,
		Age:   age,
		Email: email,
		Phone: phone,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("c.client.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return ErrNotFound

		default:
			return fmt.Errorf("request failed: status: %s body: %s", resp.Status, body)
		}
	}

	return nil
}
