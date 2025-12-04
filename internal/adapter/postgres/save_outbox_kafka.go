package postgres

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) SaveOutboxKafka(ctx context.Context, messages ...kafka.Message) error {
	if len(messages) == 0 {
		return nil
	}

	batch := make([]any, 0, len(messages))

	for _, msg := range messages {
		if msg.Topic == "" {
			return domain.ErrEmptyTopic
		}

		batch = append(batch, goqu.Record{
			"topic": msg.Topic,
			"key":   msg.Key,
			"value": msg.Value,
		})
	}

	sql, _, err := goqu.Insert("outbox").Rows(batch...).ToSQL()
	if err != nil {
		return fmt.Errorf("goqu.Insert.Rows.ToSQL: %w", err)
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err = txOrPool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
