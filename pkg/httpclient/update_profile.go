package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Update(id, name string, age int) error {
	const updateProfile = "amozhaykin/my-app/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s", c.host, updateProfile)

	request := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		ID:   id,
		Name: name,
		Age:  age,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request failed: status: %s, body: %s", resp.Status, body)
	}

	return nil
}
