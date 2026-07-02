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

func Test_UpdateProfile_Success(t *testing.T) {
	otel.SilentModeInit()
	transaction.IsUnitTest = true

	id := uuid.New()
	profile := domain.Profile{ID: id}

	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", mock.Anything, id).Return(profile, nil)
	postgres.On("UpdateProfile", mock.Anything, mock.Anything).Return(nil)
	defer postgres.AssertCalled(t, "GetProfile", mock.Anything, id)
	defer postgres.AssertCalled(t, "UpdateProfile", mock.Anything, mock.Anything)

	redis := new(mocks.Redis)
	redis.On("SetCache", mock.Anything, mock.Anything).Return(nil)
	defer redis.AssertCalled(t, "SetCache", mock.Anything, mock.Anything)

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
		require.NoError(t, err)
	}
}

func Test_UpdateProfile_InvalidUUID(t *testing.T) {
	otel.SilentModeInit()
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		name := "John Doe"

		input := dto.UpdateProfileInput{
			ID:   "invalid-uuid",
			Name: &name,
		}

		err := u.UpdateProfile(context.Background(), input)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrUUIDInvalid)
	}
}

func Test_UpdateProfile_AllFieldsAreEmpty(t *testing.T) {
	otel.SilentModeInit()
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		input := dto.UpdateProfileInput{ID: uuid.New().String()}

		err := u.UpdateProfile(context.Background(), input)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrAllFieldsAreEmpty)
	}
}

func Test_UpdateProfile_NoChanges(t *testing.T) {
	otel.SilentModeInit()
	transaction.IsUnitTest = true

	id := uuid.New()
	profile := domain.Profile{
		ID:   id,
		Name: "John Doe",
		Age:  30,
		Contacts: domain.Contacts{
			Email: "john.doe@example.com",
			Phone: "+1234567890",
		},
	}

	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", mock.Anything, id).Return(profile, nil)
	defer postgres.AssertCalled(t, "GetProfile", mock.Anything, id)

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
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrNoChangesFound)
	}
}
