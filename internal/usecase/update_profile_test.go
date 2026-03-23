package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/internal/usecase/mocks"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

// в usecase updateProfile, проверяем:
// успешный кейс
// кейс с ошибкой функции uuid.Parse (когда передан невалидный id)
// кейс с ошибкой функции input.Validate (когда все переданные для обновления поля пустые)
// кейс с ошибкой domain.ErrNoChangesFound (когда переданные поля точно такие же)
// функция profile.IsDeleted уже проверена в unit тесте на domain

func Test_UpdateProfile_Success(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	// создаём случайный UUID для теста
	id := uuid.New()
	// создаём "фиктивный" профиль, который вернёт мок базы данных при запросе Get
	profile := domain.Profile{ID: id}

	// настраиваем поведение Postgres, создаём мок (заглушку) для Postgres
	postgres := new(mocks.Postgres)
	// когда вызовут GetProfile(любой аргумент, id) → вернуть (profile, nil)
	postgres.On("GetProfile", mock.Anything, id).Return(profile, nil)
	// когда вызовут UpdateProfile(любой аргумент, любой аргумент) → вернуть (nil)
	postgres.On("UpdateProfile", mock.Anything, mock.Anything).Return(nil)
	// эти defer гарантируют, что после выполнения теста будет проверено, действительно ли методы мока вызывались
	defer postgres.AssertCalled(t, "GetProfile", mock.Anything, id)
	defer postgres.AssertCalled(t, "UpdateProfile", mock.Anything, mock.Anything)

	// Настраиваем поведение Redis (так же как и у Postgres)
	redis := new(mocks.Redis)
	redis.On("SetCache", mock.Anything, mock.Anything).Return(nil)
	defer redis.AssertCalled(t, "SetCache", mock.Anything, mock.Anything)

	// создаём экземпляр UseCase, передавая в него мок базы
	u := usecase.New(postgres, redis, nil)

	{ // сам тест
		var (
			newName  = "John Doe"
			newAge   = 30
			newEmail = "john.doe@example.com"
			newPhone = "+1234567890"
		)

		input := dto.UpdateProfileInput{
			ID:    id.String(),
			Name:  &newName,
			Age:   &newAge,
			Email: &newEmail,
			Phone: &newPhone,
		}

		err := u.UpdateProfile(context.Background(), input)
		require.NoError(t, err) // проверяем, что ошибки нет
	}
}

func Test_UpdateProfile_InvalidUUID(t *testing.T) {
	otel.SilentModeInit()
	// т.к. при невалидном ID до похода в базу дело всеравно не дойдет, то моки здесь не нужны
	// собираем UseCase
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		name := "John Doe"

		input := dto.UpdateProfileInput{
			ID:   "invalid-uuid",
			Name: &name,
		}

		err := u.UpdateProfile(context.Background(), input)
		require.Error(t, err)                          // проверяем, что ошибка действительно произошла
		require.ErrorIs(t, err, domain.ErrUUIDInvalid) // проверяем, что ошибка именно domain.ErrUUIDInvalid
	}
}

func Test_UpdateProfile_AllFieldsAreEmpty(t *testing.T) {
	otel.SilentModeInit()
	// т.к. при невалидном запросе до похода в базу дело всеравно не дойдет, то моки здесь не нужны
	// собираем UseCase
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		input := dto.UpdateProfileInput{ID: uuid.New().String()}

		err := u.UpdateProfile(context.Background(), input)
		require.Error(t, err)                                // проверяем, что ошибка действительно произошла
		require.ErrorIs(t, err, domain.ErrAllFieldsAreEmpty) // проверяем, что ошибка именно domain.ErrAllFieldsAreEmpty
	}
}

func Test_UpdateProfile_NoChanges(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	// создаём случайный UUID для теста
	id := uuid.New()
	// создаём "фиктивный" профиль, который вернёт мок базы данных при запросе Get
	profile := domain.Profile{
		ID:   id,
		Name: "John Doe",
		Age:  30,
		Contacts: domain.Contacts{
			Email: "john.doe@example.com",
			Phone: "+1234567890",
		},
	}

	// настраиваем поведение Postgres
	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", mock.Anything, id).Return(profile, nil)
	// после теста проверим, что метод реально вызывался с нужными аргументами
	defer postgres.AssertCalled(t, "GetProfile", mock.Anything, id)

	// создаём экземпляр UseCase, передавая в него мок базы
	u := usecase.New(postgres, nil, nil)

	{ // сам тест
		var (
			newName  = "John Doe"
			newAge   = 30
			newEmail = "john.doe@example.com"
			newPhone = "+1234567890"
		)

		input := dto.UpdateProfileInput{
			ID:    id.String(),
			Name:  &newName,
			Age:   &newAge,
			Email: &newEmail,
			Phone: &newPhone,
		}

		err := u.UpdateProfile(context.Background(), input)
		require.Error(t, err)                             // проверяем, что ошибка действительно произошла
		require.ErrorIs(t, err, domain.ErrNoChangesFound) // проверяем, что ошибка именно domain.ErrNoChangesFound
	}
}
