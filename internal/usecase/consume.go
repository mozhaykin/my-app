package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func (u *UseCase) Consume(ctx context.Context, msg kafka.Message) error {
	ctx, span := tracer.Start(ctx, "usecase Consume")
	defer span.End()

	if u.redis.IsIdempotencyKeyExists(ctx, string(msg.Key)) {
		log.Info().Str("key", string(msg.Key)).Msg("usecase: Consume: message already processed")

		return nil
	}

	var profile domain.Profile

	err := json.Unmarshal(msg.Value, &profile)
	if err != nil {
		return fmt.Errorf("json.Unmarshal kafka.Message: %w", err)
	}

	log.Info().
		Str("topic", msg.Topic).
		Int("partition", msg.Partition).
		Int64("offset", msg.Offset).
		Str("key", string(msg.Key)).
		Interface("profile", profile).
		Msg("consume")

	return nil
}
