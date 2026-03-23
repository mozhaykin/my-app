package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/internal/usecase/mocks"
	"github.com/mozhaykin/my-app/pkg/otel"
)

// в usecase deleteProfile, проверяем:
// успешный кейс
// кейс с ошибкой функции uuid.Parse (когда передан невалидный id)

func Test_DeleteProfile_Success(t *testing.T) {
	otel.SilentModeInit() // отключить open telemetry

	// создаём случайный UUID для теста
	id := uuid.New()

	// Создаем ошибку чтобы проверить запись в лог при отсутствии данных в Redis
	notFound := errors.New("notFound")

	// настраиваем поведение Postgres, создаём мок (заглушку) для Postgres
	postgres := new(mocks.Postgres)
	// когда вызовут DeleteProfile(любой аргумент, id) → вернуть ( nil)
	postgres.On("DeleteProfile", mock.Anything, id).Return(nil)
	// после теста проверим, что метод реально вызывался с нужными аргументами
	defer postgres.AssertCalled(t, "DeleteProfile", mock.Anything, id)

	// Настраиваем поведение Redis (так же как и у Postgres)
	redis := new(mocks.Redis)
	redis.On("DeleteCache", mock.Anything, id).Return(notFound)
	defer redis.AssertCalled(t, "DeleteCache", mock.Anything, id)

	// создаём экземпляр UseCase, передавая в него моки
	u := usecase.New(postgres, redis, nil)

	{ // сам тест
		input := dto.DeleteProfileInput{ID: id.String()}

		err := u.DeleteProfile(context.Background(), input)
		require.NoError(t, err) // проверяем, что ошибки нет
	}
}

func Test_DeleteProfile_InvalidUUID(t *testing.T) {
	otel.SilentModeInit() // отключить open telemetry
	// т.к. при невалидном ID до похода в Redis или базу дело всеравно не дойдет, то моки здесь не нужны
	// собираем UseCase
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		input := dto.DeleteProfileInput{ID: "invalid-uuid"} // строка с невалидным id

		err := u.DeleteProfile(context.Background(), input)
		require.Error(t, err)                          // проверяем, что ошибка действительно произошла
		require.ErrorIs(t, err, domain.ErrUUIDInvalid) // проверяем, что ошибка именно domain.ErrUUIDInvalid
	}
}

func Test_DeleteProfile_NotFound(t *testing.T) {
	otel.SilentModeInit() // отключить open telemetry

	// создаём случайный UUID для теста
	id := uuid.New()

	notFound := errors.New("notFound")

	// настраиваем поведение Postgres, создаём мок (заглушку) для Postgres
	postgres := new(mocks.Postgres)
	// когда вызовут DeleteProfile(любой аргумент, id) → вернуть (notFound)
	postgres.On("DeleteProfile", mock.Anything, id).Return(notFound)
	// после теста проверим, что метод реально вызывался с нужными аргументами
	defer postgres.AssertCalled(t, "DeleteProfile", mock.Anything, id)

	// создаём экземпляр UseCase, передавая в него моки
	u := usecase.New(postgres, nil, nil)

	{ // сам тест
		input := dto.DeleteProfileInput{ID: id.String()}

		err := u.DeleteProfile(context.Background(), input)
		require.ErrorIs(t, err, notFound) // проверяем, что ошибка именно notFound
	}
}
