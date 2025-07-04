package postgres

import (
	"context"
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func (p *Postgres) CreateProperty(ctx context.Context, property domain.Property) error {
	const sql = `INSERT INTO property (profile_id, tags)
                    VALUES ($1, $2)`

	args := []any{
		property.ProfileID,
		property.Tags,
	}

	tx, err := transaction.Get(ctx)
	if err != nil {
		return fmt.Errorf("transaction.Get: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("tx.Exec: %w", err)
	}

	return nil
}
