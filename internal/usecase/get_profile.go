package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) GetProfile(ctx context.Context, input dto.GetProfileInput) (dto.GetProfileOutput, error) {
	var output dto.GetProfileOutput

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return output, fmt.Errorf("uuid.Parse: %w", domain.ErrUUIDInvalid)
	}

	// Запрос в Redis
	profile, err := u.redis.GetCache(ctx, id)

	// Если данные получены и ошибки нет, то выходим
	if err == nil {
		output.Profile = profile

		return output, nil
	}

	// Если ошибка не ErrNotFound, то пишем ее в лог
	if !errors.Is(err, domain.ErrNotFound) {
		log.Error().Err(err).Str("profileID", id.String()).Msg("cache: GetProfile: get cache")
	}

	// Выполняем запрос к базе данных
	profile, err = u.postgres.GetProfile(ctx, id)
	if err != nil {
		return output, fmt.Errorf("u.postgres.GetProfile: %w", err)
	}

	if profile.IsDeleted() {
		return output, fmt.Errorf("profile.IsDeleted: %w", domain.ErrNotFound)
	}

	// Set в Redis (после получения данных из базы)
	err = u.redis.SetCache(ctx, profile)
	if err != nil {
		log.Error().Err(err).Str("profileID", id.String()).Msg("cache: GetProfile: set cache")
	}

	output.Profile = profile

	return output, nil
}
