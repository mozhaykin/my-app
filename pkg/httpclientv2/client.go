package httpclientv2

import (
	"errors"
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_client"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	client *http_client.ClientWithResponses
}

func New(host string) (*Client, error) {
	client, err := http_client.NewClientWithResponses(host)
	if err != nil {
		return nil, fmt.Errorf("http_client.NewClient: %w", err)
	}

	return &Client{client: client}, nil
}
