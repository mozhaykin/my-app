package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

func (p *Postgres) GetProfile(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	const sql = `SELECT created_at, updated_at, deleted_at, name, age, status, verified, contacts 
				FROM profile WHERE id = $1`

	dto := struct {
		CreatedAt pgtype.Timestamptz
		UpdatedAt pgtype.Timestamptz
		DeletedAt pgtype.Timestamptz
		Name      pgtype.Text
		Age       pgtype.Int4
		Status    pgtype.Text
		Verified  pgtype.Bool
		Contacts  []byte
	}{}

	dest := []any{
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&dto.DeletedAt,
		&dto.Name,
		&dto.Age,
		&dto.Status,
		&dto.Verified,
		&dto.Contacts,
	}

	err := p.pool.QueryRow(ctx, sql, profileID).Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Profile{}, fmt.Errorf("p.pool.QueryRow.Scan: %w", domain.ErrNotFound)
		}

		return domain.Profile{}, fmt.Errorf("p.pool.QueryRow: %w", err)
	}

	var contacts domain.Contacts

	err = json.Unmarshal(dto.Contacts, &contacts)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	profile := domain.Profile{
		ID:        profileID,
		CreatedAt: dto.CreatedAt.Time,
		UpdatedAt: dto.UpdatedAt.Time,
		DeletedAt: dto.DeletedAt.Time,
		Name:      domain.Name(dto.Name.String),
		Age:       domain.Age(dto.Age.Int32),
		Status:    domain.NewStatus(dto.Status.String),
		Verified:  dto.Verified.Bool,
		Contacts:  contacts,
	}

	return profile, nil
}
