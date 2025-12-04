package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

//go:generate mockery

type Postgres interface {
	CreateProfile(ctx context.Context, profile domain.Profile) error
	CreateProperty(ctx context.Context, property domain.Property) error
	GetProfile(ctx context.Context, profileID uuid.UUID) (domain.Profile, error)
	GetProfiles(ctx context.Context, input dto.GetProfilesInput) ([]domain.Profile, error)
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	DeleteProfile(ctx context.Context, profileID uuid.UUID) error

	SaveOutboxKafka(ctx context.Context, events ...kafka.Message) error
	ReadOutboxKafka(ctx context.Context, limit int) ([]kafka.Message, error)
}

type Redis interface {
	GetCache(ctx context.Context, id uuid.UUID) (domain.Profile, error)
	SetCache(ctx context.Context, profile domain.Profile) error
	DeleteCache(ctx context.Context, id uuid.UUID) error
}

type Kafka interface {
	Produce(ctx context.Context, events ...kafka.Message) error
}

type UseCase struct {
	postgres Postgres
	redis    Redis
	kafka    Kafka
}

func New(postgres Postgres, redis Redis, k Kafka) *UseCase {
	return &UseCase{
		postgres: postgres,
		redis:    redis,
		kafka:    k,
	}
}
