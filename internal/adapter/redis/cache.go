package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func (r *Redis) GetCache(ctx context.Context, id uuid.UUID) (domain.Profile, error) {
	ctx, span := tracer.Start(ctx, "adapter rediscache GetCache")
	defer span.End()

	var profile domain.Profile

	key := cacheKeyPrefix + id.String()

	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return profile, domain.ErrNotFound
		}

		return profile, fmt.Errorf("r.client.Get: %w", err)
	}

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return profile, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return profile, nil
}

func (r *Redis) SetCache(ctx context.Context, profile domain.Profile) error {
	ctx, span := tracer.Start(ctx, "adapter rediscache SetCache")
	defer span.End()

	data, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	key := cacheKeyPrefix + profile.ID.String()

	err = r.redis.Set(ctx, key, data, cacheTTL).Err()
	if err != nil {
		return fmt.Errorf("r.client.Set: %w", err)
	}

	return nil
}

func (r *Redis) DeleteCache(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "adapter rediscache DeleteCache")
	defer span.End()

	key := cacheKeyPrefix + id.String()

	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("r.client.Del: %w", err)
	}

	return nil
}
