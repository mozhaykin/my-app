package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      Name      `json:"name"`
	Age       Age       `json:"age"`
	Status    Status    `json:"status"`
	Verified  bool      `json:"verified"`
	Contacts  Contacts  `json:"contacts"`
}
type Name string

type Age int

type Status int

type Contacts struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (c *Client) Get(ctx context.Context, id string) (Profile, error) {
	const getProfile = "amozhaykin/my-app/api/v1/profile"

	path := fmt.Sprintf("http://%s/%s/%s", c.host, getProfile, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return Profile{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close() //nolint:contextcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		switch {
		case resp.StatusCode == http.StatusNotFound:
			return Profile{}, ErrNotFound

		default:
			return Profile{}, fmt.Errorf("request failed: status: %s body: %s", resp.Status, body)
		}
	}

	var profile Profile

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return Profile{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return profile, nil
}
