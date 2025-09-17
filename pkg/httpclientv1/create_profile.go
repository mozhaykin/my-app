package httpclientv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type CreateProfileRequest struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (c *Client) Create(request CreateProfileRequest) (uuid.UUID, error) {
	const createProfile = "amozhaykin/my-app/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s", c.host, createProfile)

	body, err := json.Marshal(request)
	if err != nil {
		return uuid.Nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return uuid.Nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return uuid.Nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return uuid.Nil, fmt.Errorf("request failed: status: %s, body: %s", resp.Status, body)
	}

	response := struct {
		ID uuid.UUID `json:"id"`
	}{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return uuid.Nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return response.ID, nil
}
