package usecase

import (
	"fmt"

	"golang.org/x/net/context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

//nolint:nonamedreturns
func (u *UseCase) OutboxReadAndProduce(ctx context.Context, limit int) (count int, err error) {
	ctx, err = transaction.Begin(ctx)
	if err != nil {
		return count, fmt.Errorf("transaction.Begin: %w", err)
	}

	defer transaction.Rollback(ctx)

	// Читаем сообщения из outbox таблицы БД
	msgs, err := u.postgres.ReadOutboxKafka(ctx, limit)
	if err != nil {
		return count, fmt.Errorf("u.postgres.ReadOutboxKafka: %w", err)
	}

	count = len(msgs)

	// Пишем в Kafka
	err = u.kafka.Produce(ctx, msgs...)
	if err != nil {
		return count, fmt.Errorf("u.kafka.Produce: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return count, fmt.Errorf("transaction.Commit: %w", err)
	}

	return count, nil
}
