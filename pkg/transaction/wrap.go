package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func Wrap(ctx context.Context, fn func(context.Context) error) error {
	if IsUnitTest {
		return fn(ctx)
	}

	ctx, span := tracer.Start(ctx, "transaction Wrap")
	defer span.End()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Msg("transaction: Rollback")
		}
	}()

	// Создаем дочерний контекст с ключем (пустая структура) и значением Transaction
	ctx = context.WithValue(ctx, ctxKey{}, &Transaction{tx})

	// Передаем дочерний контекст в анонимную функцию, чтобы потом извлечь его в функции TryExtractTX
	err = fn(ctx)
	if err != nil {
		return fmt.Errorf("fn: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
