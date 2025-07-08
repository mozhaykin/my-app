package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(p *postgres.Pool) *Postgres {
	return &Postgres{
		pool: p.Pool,
	}
}
