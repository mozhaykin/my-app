package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func (u *UseCase) OutboxReadAndProduce(ctx context.Context, limit int) (int, error) {
	// Создаем новый трейс, указываем spanName(название пакета и функция)
	ctx, span := tracer.Start(ctx, "usecase OutboxReadAndProduce")
	defer span.End() // Обязательно закрываем span

	var events []domain.Event

	// Транзакция на чтение из outbox
	err := transaction.Wrap(ctx, func(ctx context.Context) error {
		// Читаем событие из outbox таблицы БД
		var err error

		events, err = u.postgres.ReadOutbox(ctx, limit)
		if err != nil {
			return fmt.Errorf("u.postgres.ReadOutbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("transaction(read): %w", err)
	}

	eventsCount := len(events)
	// Если events нет, выходим из функции без ошибок
	if eventsCount == 0 {
		return 0, nil
	}

	// Вне транзакции публикуем события в kafka
	err = u.kafka.Produce(ctx, events)
	if err != nil {
		return eventsCount, fmt.Errorf("u.kafka.Produce: %w", err)
	}

	// Формируем список id для удаления events из outbox
	ids := extractEventIDs(events)

	// Новая транзакция для удаления events из outbox
	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		err = u.postgres.ClearOutbox(ctx, ids)
		if err != nil {
			return fmt.Errorf("u.postgres.ClearOutbox: %w", err)
		}

		return nil
	})
	if err != nil {
		return eventsCount, fmt.Errorf("transaction(clear): %w", err)
	}

	log.Info().Int("msgs", eventsCount).Msg("outbox events produce")

	return eventsCount, nil
}

func extractEventIDs(events []domain.Event) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.ID)
	}

	return ids
}
