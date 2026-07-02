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

func Test_DeleteProfile_Success(t *testing.T) {
	otel.SilentModeInit()

	id := uuid.New()

	notFound := errors.New("notFound")

	postgres := new(mocks.Postgres)
	postgres.On("DeleteProfile", mock.Anything, id).Return(nil)
	defer postgres.AssertCalled(t, "DeleteProfile", mock.Anything, id)

	redis := new(mocks.Redis)
	redis.On("DeleteCache", mock.Anything, id).Return(notFound)
	defer redis.AssertCalled(t, "DeleteCache", mock.Anything, id)

	u := usecase.New(postgres, redis, nil)

	{ // сам тест
		input := dto.DeleteProfileInput{ID: id.String()}

		err := u.DeleteProfile(context.Background(), input)
		require.NoError(t, err)
	}
}

func Test_DeleteProfile_InvalidUUID(t *testing.T) {
	otel.SilentModeInit()
	u := usecase.New(nil, nil, nil)

	{ // Сам тест
		input := dto.DeleteProfileInput{ID: "invalid-uuid"}

		err := u.DeleteProfile(context.Background(), input)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrUUIDInvalid)
	}
}

func Test_DeleteProfile_NotFound(t *testing.T) {
	otel.SilentModeInit()

	id := uuid.New()

	notFound := errors.New("notFound")

	postgres := new(mocks.Postgres)
	postgres.On("DeleteProfile", mock.Anything, id).Return(notFound)
	defer postgres.AssertCalled(t, "DeleteProfile", mock.Anything, id)

	u := usecase.New(postgres, nil, nil)

	{ // сам тест
		input := dto.DeleteProfileInput{ID: id.String()}

		err := u.DeleteProfile(context.Background(), input)
		require.ErrorIs(t, err, notFound)
	}
}
