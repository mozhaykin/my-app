package usecase_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase/mocks"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func Test_Consume_Success(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	ctx := context.Background()

	redis := new(mocks.Redis)
	u := usecase.New(nil, redis, nil)

	profile := domain.Profile{
		ID:   uuid.New(),
		Name: "John",
	}
	data, err := json.Marshal(profile)
	require.NoError(t, err)

	msg := kafka.Message{
		Key:   []byte("key-success"),
		Value: data,
	}

	redis.On("IsIdempotencyKeyExists", mock.Anything, "key-success").Once().Return(false)

	err = u.Consume(ctx, msg)
	require.NoError(t, err)

	redis.AssertExpectations(t)
}

func Test_Consume_AlreadyProcessed(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	ctx := context.Background()

	redis := new(mocks.Redis)
	u := usecase.New(nil, redis, nil)

	msg := kafka.Message{
		Key: []byte("key-processed"),
	}

	redis.On("IsIdempotencyKeyExists", mock.Anything, "key-processed").Once().Return(true)

	err := u.Consume(ctx, msg)
	require.NoError(t, err)

	redis.AssertExpectations(t)
}

func Test_Consume_InvalidJSON(t *testing.T) {
	otel.SilentModeInit()         // отключить open telemetry
	transaction.IsUnitTest = true // отключить транзакции

	ctx := context.Background()

	redis := new(mocks.Redis)
	u := usecase.New(nil, redis, nil)

	msg := kafka.Message{
		Key:   []byte("key-invalid-json"),
		Value: []byte("invalid-json"),
	}

	redis.On("IsIdempotencyKeyExists", mock.Anything, "key-invalid-json").Once().Return(false)

	err := u.Consume(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "json.Unmarshal")

	redis.AssertExpectations(t)
}
