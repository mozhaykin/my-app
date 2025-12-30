package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase/mocks"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

// в usecase createProfile проверяем только успешный кейс
// функция domain.NewProfile уже проверена в unit тесте на domain
// функцию domain.NewProperty проверять нет смысла т.к. там просто присвоение полей

func Test_CreateProfile_Success(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	// настраиваем поведение Postgres, создаём мок (заглушку) для Postgres
	postgres := new(mocks.Postgres)
	// когда вызовут CreateProfile(любой аргумент, любой аргумент) → вернуть (nil)
	postgres.On("CreateProfile", mock.Anything, mock.Anything).Return(nil)
	// когда вызовут CreateProperty(любой аргумент, любой аргумент) → вернуть (nil)
	postgres.On("CreateProperty", mock.Anything, mock.Anything).Return(nil)
	// когда вызовут SaveOutboxKafka(любой аргумент, любой аргумент) → вернуть (nil)
	postgres.On("SaveOutbox", mock.Anything, mock.Anything).Return(nil)
	// эти defer гарантируют, что после выполнения теста будет проверено, действительно ли методы мока вызывались
	defer postgres.AssertCalled(t, "CreateProfile", mock.Anything, mock.Anything)
	defer postgres.AssertCalled(t, "CreateProperty", mock.Anything, mock.Anything)
	defer postgres.AssertCalled(t, "SaveOutbox", mock.Anything, mock.Anything)

	// создаём экземпляр UseCase, передавая в него мок базы
	u := usecase.New(postgres, nil, nil)

	{ // сам тест
		input := dto.CreateProfileInput{
			Name:  "John Doe",
			Age:   30,
			Email: "john.doe@example.com",
			Phone: "+1234567890",
		}

		actual, err := u.CreateProfile(context.Background(), input)
		require.NoError(t, err)        // проверяем, что ошибки нет
		require.NotEmpty(t, actual.ID) // проверяем, что в ответе вернулся непустой ID профиля
	}
}
