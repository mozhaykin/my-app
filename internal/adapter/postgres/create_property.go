package postgres

import (
	"context"
	"fmt"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) CreateProperty(ctx context.Context, property domain.Property) error {
	ctx, span := tracer.Start(ctx, "adapter postgres CreateProperty")
	defer span.End()

	const sql = `INSERT INTO property (profile_id, tags) VALUES ($1, $2)`

	args := []any{
		property.ProfileID,
		property.Tags,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err := txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
