package httpclientv2

import (
	"errors"
	"fmt"

	"github.com/mozhaykin/my-app/gen/http/profile_v2/client"
)

var ErrNotFound = errors.New("not found")

type Config struct {
	Address string `envconfig:"HTTP_CLIENT_V2_ADDRESS" required:"true"`
}

type Client struct {
	client *client.ClientWithResponses
}

func New(c Config) (*Client, error) {
	newClient, err := client.NewClientWithResponses(c.Address)
	if err != nil {
		return nil, fmt.Errorf("http_client.NewClient: %w", err)
	}

	return &Client{client: newClient}, nil
}
