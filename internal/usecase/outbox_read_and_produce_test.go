package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase/mocks"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func Test_OutboxReadAndProduce_Success(t *testing.T) {
	transaction.IsUnitTest = true

	// Данные для поведения
	msgs := []domain.Event{{}, {}, {}}

	// Настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("ReadOutboxKafka", mock.Anything, mock.Anything).Return(msgs, nil)
	defer postgres.AssertCalled(t, "ReadOutboxKafka", mock.Anything, mock.Anything)

	// Настройка поведения kafka
	kafka := new(mocks.Kafka)
	kafka.On("Produce", mock.Anything, mock.Anything).Return(nil)
	defer kafka.AssertCalled(t, "Produce", mock.Anything, mock.Anything)

	// Создаём экземпляр UseCase, передавая в него мок базы
	u := usecase.New(postgres, nil, kafka)

	{ // Сам тест
		count, err := u.OutboxReadAndProduce(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, len(msgs), count)
	}
}
