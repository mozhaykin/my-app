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
	ctx, span := tracer.Start(ctx, "usecase GetProfile")
	defer span.End()

	var output dto.GetProfileOutput

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return output, fmt.Errorf("uuid.Parse: %w", domain.ErrUUIDInvalid)
	}

	profile, err := u.redis.GetCache(ctx, id)

	if err == nil {
		output.Profile = profile

		return output, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		log.Error().Err(err).Str("profileID", id.String()).Msg("cache: GetProfile: get cache")
	}

	profile, err = u.postgres.GetProfile(ctx, id)
	if err != nil {
		return output, fmt.Errorf("u.postgres.GetProfile: %w", err)
	}

	if profile.IsDeleted() {
		return output, fmt.Errorf("profile.IsDeleted: %w", domain.ErrNotFound)
	}

	err = u.redis.SetCache(ctx, profile)
	if err != nil {
		log.Error().Err(err).Str("profileID", id.String()).Msg("cache: GetProfile: set cache")
	}

	output.Profile = profile

	return output, nil
}
