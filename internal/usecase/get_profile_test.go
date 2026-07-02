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

func Test_GetProfile_Success(t *testing.T) {
	otel.SilentModeInit()

	id := uuid.New()
	profile := domain.Profile{ID: id}

	someError := errors.New("SomeError")

	redis := new(mocks.Redis)
	redis.On("GetCache", mock.Anything, id).Return(profile, someError)
	redis.On("SetCache", mock.Anything, profile).Return(nil)
	defer redis.AssertCalled(t, "GetCache", mock.Anything, id)
	defer redis.AssertCalled(t, "SetCache", mock.Anything, profile)

	postgres := new(mocks.Postgres)
	postgres.On("GetProfile", mock.Anything, id).Return(profile, nil)
	defer postgres.AssertCalled(t, "GetProfile", mock.Anything, id)

	u := usecase.New(postgres, redis, nil)

	{ // сам тест
		input := dto.GetProfileInput{ID: id.String()}
		output := dto.GetProfileOutput{Profile: profile}

		actual, err := u.GetProfile(context.Background(), input)
		require.NoError(t, err)
		require.Equal(t, output, actual)
	}
}

func Test_GetProfile_InvalidUUID(t *testing.T) {
	otel.SilentModeInit()
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		input := dto.GetProfileInput{ID: "invalid-uuid"}
		output := dto.GetProfileOutput{}

		actual, err := u.GetProfile(context.Background(), input)
		require.Error(t, err)
		require.Equal(t, output, actual)
		require.ErrorIs(t, err, domain.ErrUUIDInvalid)
	}
}
