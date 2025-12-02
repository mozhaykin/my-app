package usecase

import (
	"context"

	"github.com/google/uuid"

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

	SaveOutboxKafka(ctx context.Context, events ...domain.Event) error
	ReadOutboxKafka(ctx context.Context, limit int) ([]domain.Event, error)
}

type Redis interface {
	GetCache(context.Context, uuid.UUID) (domain.Profile, error)
	SetCache(context.Context, domain.Profile) error
	DeleteCache(context.Context, uuid.UUID) error
}

type Kafka interface {
	Produce(ctx context.Context, events ...domain.Event) error
}

type UseCase struct {
	postgres Postgres
	redis    Redis
	kafka    Kafka
}

func New(postgres Postgres, redis Redis, kafka Kafka) *UseCase {
	return &UseCase{
		postgres: postgres,
		redis:    redis,
		kafka:    kafka,
	}
}
