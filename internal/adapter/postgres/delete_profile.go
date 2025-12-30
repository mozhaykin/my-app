package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel/tracer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) DeleteProfile(ctx context.Context, profileID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "adapter postgres DeleteProfile")
	defer span.End()

	const sql = `UPDATE profile SET deleted_at = NOW() 
               WHERE id = $1`

	txOrPool := transaction.TryExtractTX(ctx)

	result, err := txOrPool.Exec(ctx, sql, profileID)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("result.RowsAffected: %w", domain.ErrNotFound)
	}

	return nil
}
