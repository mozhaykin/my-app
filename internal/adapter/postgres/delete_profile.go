package postgres

import (
	"context"
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"

	"github.com/google/uuid"
)

func (p *Postgres) DeleteProfile(ctx context.Context, profileID uuid.UUID) error {
	const sql = `UPDATE profile SET deleted_at = NOW() WHERE id = $1`

	result, err := p.pool.Exec(ctx, sql, profileID)
	if err != nil {
		return fmt.Errorf("p.pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("result.RowsAffected: %w", domain.ErrNotFound)
	}

	return nil
}
