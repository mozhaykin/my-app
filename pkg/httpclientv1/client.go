package httpclientv1

import (
	"errors"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	client http.Client
	host   string
}

func New(host string) *Client {
	return &Client{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		host: host,
	}
}
