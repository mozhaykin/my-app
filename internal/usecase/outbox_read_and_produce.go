package usecase

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

//nolint:nonamedreturns
func (u *UseCase) OutboxReadAndProduce(ctx context.Context, limit int) (batchSize int, err error) {
	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		// Читаем сообщения из outbox таблицы БД
		batch, err := u.postgres.ReadOutboxKafka(ctx, limit)
		if err != nil {
			return fmt.Errorf("u.postgres.ReadOutboxKafka: %w", err)
		}

		batchSize = len(batch)

		// Пишем в Kafka
		err = u.kafka.Produce(ctx, batch...)
		if err != nil {
			return fmt.Errorf("u.kafka.Produce: %w", err)
		}

		return nil
	})
	if err != nil {
		return batchSize, fmt.Errorf("transaction.Wrap: %w", err)
	}

	log.Info().Int("msgs", batchSize).Msg("outbox kafka read and produce")

	return batchSize, nil
}
