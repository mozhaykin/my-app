package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mozhaykin/my-app/pkg/postgres"
)

//nolint:gochecknoglobals
var (
	pool       *pgxpool.Pool
	IsUnitTest bool
)

func Init(p *postgres.Pool) {
	pool = p.Pool
}

type ctxKey struct{}

type Transaction struct {
	pgx.Tx
}

type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func ExtractTX(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(ctxKey{}).(pgx.Tx)

	return tx, ok
}

func TryExtractTX(ctx context.Context) Executor {
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return pool
	}

	return tx
}
