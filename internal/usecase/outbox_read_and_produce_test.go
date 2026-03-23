package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/internal/usecase/mocks"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func Test_OutboxReadAndProduce_Success(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	ctx := context.Background()
	limit := 10

	// Данные для поведения
	events := []domain.Event{
		{ID: uuid.New()},
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	expectedIDs := []uuid.UUID{
		events[0].ID,
		events[1].ID,
		events[2].ID,
	}

	// Настраиваем поведение Postgres и kafka
	postgres := new(mocks.Postgres)
	k := new(mocks.Kafka)

	// Ожидаемый порядок вызовов
	mock.InOrder(
		postgres.On("ReadOutbox", mock.Anything, limit).Once().Return(events, nil),
		k.On("Produce", mock.Anything, events).Once().Return(nil),
		postgres.On("ClearOutbox", mock.Anything, expectedIDs).Once().Return(nil),
	)

	defer postgres.AssertExpectations(t)
	defer k.AssertExpectations(t)

	// Создаём экземпляр UseCase, передавая в него мок базы
	u := usecase.New(postgres, nil, k)

	{ // Сам тест
		count, err := u.OutboxReadAndProduce(ctx, limit)
		require.NoError(t, err)
		require.Equal(t, len(events), count)
	}
}
