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

	SaveOutbox(ctx context.Context, events ...domain.Event) error
	ReadOutbox(ctx context.Context, limit int) ([]domain.Event, error)
	ClearOutbox(ctx context.Context, ids []uuid.UUID) error
}

type Redis interface {
	GetCache(ctx context.Context, id uuid.UUID) (domain.Profile, error)
	SetCache(ctx context.Context, profile domain.Profile) error
	DeleteCache(ctx context.Context, id uuid.UUID) error

	IsIdempotencyKeyExists(ctx context.Context, idempotencyKey string) bool
}

type Kafka interface {
	Produce(ctx context.Context, events []domain.Event) error
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
