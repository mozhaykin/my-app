package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/internal/usecase/mocks"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func Test_CreateProfile_Success(t *testing.T) {
	otel.SilentModeInit()
	transaction.IsUnitTest = true

	postgres := new(mocks.Postgres)
	postgres.On("CreateProfile", mock.Anything, mock.Anything).Return(nil)
	postgres.On("CreateProperty", mock.Anything, mock.Anything).Return(nil)
	postgres.On("SaveOutbox", mock.Anything, mock.Anything).Return(nil)
	defer postgres.AssertCalled(t, "CreateProfile", mock.Anything, mock.Anything)
	defer postgres.AssertCalled(t, "CreateProperty", mock.Anything, mock.Anything)
	defer postgres.AssertCalled(t, "SaveOutbox", mock.Anything, mock.Anything)

	u := usecase.New(postgres, nil, nil)

	{ // сам тест
		input := dto.CreateProfileInput{
			Name:  "John Doe",
			Age:   30,
			Email: "john.doe@example.com",
			Phone: "+1234567890",
		}

		actual, err := u.CreateProfile(context.Background(), input)
		require.NoError(t, err)
		require.NotEmpty(t, actual.ID)
	}
}
