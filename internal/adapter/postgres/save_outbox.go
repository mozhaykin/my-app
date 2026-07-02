package postgres

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) SaveOutbox(ctx context.Context, events ...domain.Event) error {
	ctx, span := tracer.Start(ctx, "adapter postgres SaveOutbox")
	defer span.End()

	batch := make([]any, 0, len(events))

	for _, e := range events {
		batch = append(batch, goqu.Record{
			"event_id":      e.ID,
			"event_type":    e.Type,
			"occurred_at":   e.OccurredAt,
			"value":         e.Value,
			"trace_context": e.TraceContext, // jsonb
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
