package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) DeleteProfile(ctx context.Context, input dto.DeleteProfileInput) error {
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
