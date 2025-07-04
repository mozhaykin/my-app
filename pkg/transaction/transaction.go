package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey struct{}

var (
	errMissingInit  = errors.New("missing `transaction.Init' call before `transaction.Begin'")
	errMissingBegin = errors.New("missing `transaction.Begin' call before 'transaction.Get'")
)

//nolint:gochecknoglobals
var (
	pool       *pgxpool.Pool
	IsUnitTest bool
)

func Init(p *postgres.Pool) {
	pool = p.Pool
}

type Transaction struct {
	pgx.Tx
}

func Begin(ctx context.Context) (context.Context, error) {
	if IsUnitTest {
		return ctx, nil
	}

	if pool == nil {
		return nil, errMissingInit
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Begin: %w", err)
	}

	ctx = context.WithValue(ctx, ctxKey{}, &Transaction{tx})

	return ctx, nil
}

func Rollback(ctx context.Context) {
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return
	}

	err := tx.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		log.Error().Err(err).Msg("transaction: Rollback")
	}
}

func Commit(ctx context.Context) error {
	if IsUnitTest {
		return nil
	}

	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return errMissingBegin
	}

	err := tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

func Get(ctx context.Context) (*Transaction, error) {
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return nil, errMissingBegin
	}

	return tx, nil
}
