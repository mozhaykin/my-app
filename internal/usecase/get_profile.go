package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func (u *UseCase) GetProfile(ctx context.Context, input dto.GetProfileInput) (dto.GetProfileOutput, error) {
	// Создаем новый трейс, указываем spanName(название пакета и функция)
	ctx, span := tracer.Start(ctx, "usecase GetProfile")
	defer span.End() // Обязательно закрываем span

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
