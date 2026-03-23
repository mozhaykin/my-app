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

	// Если мы вдруг получим сообщение из kafka повторно, то нужно сделать так, чтобы его повторная обработка
	// ничего не меняла. То есть операция была идемпотентной.
	// Для этого я передаю в redis idempotencyKey и проверяю, записан ли он уже в redis.
	// Если он уже есть, значит это сообщение мы уже получали и повторно его обрабатывать не нужно, выходим.
	// Если idempotencyKey в redis нет, то я записываю его туда и обрабатываю сообщение т.к. оно получено впервые.
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
