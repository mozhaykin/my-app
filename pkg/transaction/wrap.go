package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel/tracer"
)

// Функция - обертка для транзакций
// Вызываем эту функцию и передаем в нее ctx и анонимную функцию, в которой написали код,
// который должен быть выполнен в транзакции.

func Wrap(ctx context.Context, fn func(context.Context) error) error {
	// Если IsUnitTest bool, то сразу выполняем анонимную функцию и выходим
	if IsUnitTest {
		return fn(ctx)
	}

	// Создаем новый трейс, указываем spanName(название пакета и функция)
	ctx, span := tracer.Start(ctx, "transaction Wrap")
	defer span.End() // Обязательно закрываем span

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
