package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/mozhaykin/my-app/pkg/otel/tracer"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) ClearOutbox(ctx context.Context, ids []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "adapter postgres DeleteOutbox")
	defer span.End()

	tx, ok := transaction.ExtractTX(ctx)
	if !ok {
		return errors.New("ClearOutbox must be called inside transaction")
	}

	const sql = `DELETE FROM outbox
				WHERE event_id = ANY($1)`

	_, err := tx.Exec(ctx, sql, ids)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}
