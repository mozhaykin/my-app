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

type Kafka interface {
	Produce(ctx context.Context, events ...domain.Event) error
}

type UseCase struct {
	postgres Postgres
	kafka    Kafka
}

func New(postgres Postgres, kafka Kafka) *UseCase {
	return &UseCase{
		postgres: postgres,
		kafka:    kafka,
	}
}
