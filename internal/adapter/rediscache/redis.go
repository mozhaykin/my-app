package rediscache

import (
	"time"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/redisclient"
)

const (
	prefix = "amozhaykin:my-app:"
	ttl    = 10 * time.Minute
)

type Redis struct {
	redis *redisclient.Client
}

func New(client *redisclient.Client) *Redis {
	return &Redis{
		redis: client,
	}
}
