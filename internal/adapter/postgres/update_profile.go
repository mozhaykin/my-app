package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

func (p *Postgres) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	const sql = `UPDATE profile SET name = $1, age = $2, contacts = $3, updated_at = NOW()
                     WHERE id = $4`

	contacts, err := json.Marshal(profile.Contacts)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	args := []any{
		profile.Name,
		profile.Age,
		contacts,
		profile.ID,
	}

	_, err = p.pool.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}

		return fmt.Errorf("p.pool.Exec: %w", err)
	}

	return nil
}
