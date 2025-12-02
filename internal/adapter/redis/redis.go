package redis

import (
	"time"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/redis"
)

const (
	prefix = "amozhaykin:my-app:"
	ttl    = time.Minute
)

type Redis struct {
	redis *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{
		redis: client,
	}
}
