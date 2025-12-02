package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase/mocks"
)

// в usecase getProfile, проверяем:
// успешный кейс
// кейс с ошибкой функции uuid.Parse (когда передан невалидный id)
// функция profile.IsDeleted уже проверена в unit тесте на domain

func Test_GetProfile_Success(t *testing.T) {
	// otel.SilentModeInit() // отключить open telemetry

	// создаём случайный UUID для теста
	id := uuid.New()
	// создаём "фиктивный" профиль, который вернёт мок базы данных
	profile := domain.Profile{ID: id}

	// Создаем ошибку которую вернет Redis, для того чтобы мы не вышли из основной функции раньше времени
	// и выполнение кода продолжилось
	SomeError := errors.New("SomeError")

	// Настраиваем поведение Redis, создаём мок (заглушку)
	redis := new(mocks.Redis)
	redis.On("GetCache", mock.Anything, id).Return(profile, SomeError)
	redis.On("SetCache", mock.Anything, profile).Return(nil)
	defer redis.AssertCalled(t, "GetCache", mock.Anything, id)
	defer redis.AssertCalled(t, "SetCache", mock.Anything, profile)

	// Настраиваем поведение Postgres, создаём мок (заглушку)
	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", mock.Anything, id).Return(profile, nil)
	defer postgres.AssertCalled(t, "GetProfile", mock.Anything, id)

	// создаём экземпляр UseCase, передавая в него моки
	u := usecase.New(postgres, redis, nil)

	{ // сам тест
		input := dto.GetProfileInput{ID: id.String()}
		output := dto.GetProfileOutput{Profile: profile}

		actual, err := u.GetProfile(context.Background(), input)
		require.NoError(t, err)          // проверяем, что ошибки нет
		require.Equal(t, output, actual) // сравниваем на равенство
	}
}

func Test_GetProfile_InvalidUUID(t *testing.T) {
	// otel.SilentModeInit()
	// т.к. при невалидном ID до похода в Redis или базу дело всеравно не дойдет, то моки здесь не нужны
	// Собираем UseCase
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		input := dto.GetProfileInput{ID: "invalid-uuid"} // невалидное значение id
		output := dto.GetProfileOutput{}                 // ожидаемый результат: пустой output

		actual, err := u.GetProfile(context.Background(), input)
		require.Error(t, err)                          // проверяем, что ошибка действительно произошла
		require.Equal(t, output, actual)               // проверяем, что озвратился пустой output
		require.ErrorIs(t, err, domain.ErrUUIDInvalid) // проверяем, что ошибка именно domain.ErrUUIDInvalid
	}
}
