package redis

import (
	"time"

	"github.com/mozhaykin/my-app/pkg/redisclient"
)

const (
	redisNamespace       = "amozhaykin:my-app:"
	cacheKeyPrefix       = redisNamespace + "cache:"
	idempotencyKeyPrefix = redisNamespace + "idempotency:"
	cacheTTL             = 10 * time.Minute
	idempotencyTTL       = 6 * time.Hour
)

type Redis struct {
	redis *redisclient.Client
}

func New(client *redisclient.Client) *Redis {
	return &Redis{
		redis: client,
	}
}
