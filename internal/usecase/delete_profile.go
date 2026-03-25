package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func (u *UseCase) DeleteProfile(ctx context.Context, input dto.DeleteProfileInput) error {
	ctx, span := tracer.Start(ctx, "usecase DeleteProfile")
	defer span.End()

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", domain.ErrUUIDInvalid)
	}

	err = u.postgres.DeleteProfile(ctx, id)
	if err != nil {
		return fmt.Errorf("u.postgres.DeleteProfile: %w", err)
	}

	err = u.redis.DeleteCache(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("profileID", id.String()).Msg("cache: DeleteProfile: delete cache")
	}

	return nil
}
