package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) ReadOutbox(ctx context.Context, limit int) ([]domain.Event, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres ReadOutbox")
	defer span.End()

	tx, ok := transaction.ExtractTX(ctx)
	if !ok {
		return nil, errors.New("ReadOutbox must be called inside transaction")
	}

	const sql = `SELECT event_id, event_type, occurred_at, value, trace_context
				FROM outbox
				ORDER BY created_at
				LIMIT $1
				FOR UPDATE SKIP LOCKED;`

	rows, err := tx.Query(ctx, sql, limit)
	if err != nil {
		return nil, fmt.Errorf("tx.Query: %w", err)
	}
	defer rows.Close()

	var events []domain.Event

	for rows.Next() {
		var e domain.Event

		err = rows.Scan(
			&e.ID,
			&e.Type,
			&e.OccurredAt,
			&e.Value,
			&e.TraceContext,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		events = append(events, e)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return events, nil
}
