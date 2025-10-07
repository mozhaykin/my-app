package httpclientv2

import (
	"errors"
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/client"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	client *client.ClientWithResponses
}

func New(host string) (*Client, error) {
	c, err := client.NewClientWithResponses(host)
	if err != nil {
		return nil, fmt.Errorf("http_client.NewClient: %w", err)
	}

	return &Client{client: c}, nil
}
