package transaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
)

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

type ctxKey struct{}

type Transaction struct {
	pgx.Tx
}

// Создаем интерфейс Executor, чтобы мы могли вернуть из функции TryExtractTX
// хоть pgx.Tx, хоть *pgxpool.Pool, т.к. у них обоих есть эти нужные нам методы

type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Эта функция позволяет нам писать одинаковый код при запросах к базе данных в адаптерах,
// независимо от того в транзакции мы хотим выполнять запрос или без.

func TryExtractTX(ctx context.Context) Executor {
	// Пробуем извлечь транзакцию
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	// Если в ctx нет ключа с транзакцией, то возвращаем *pgxpool.Pool, чтобы запрос к базе был без транзакции
	if !ok {
		return pool
	}

	// Если ключ есть, то возвращаем pgx.Tx, чтобы запрос к базе был в транзакции
	return tx
}
