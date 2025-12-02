package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func (u *UseCase) UpdateProfile(ctx context.Context, input dto.UpdateProfileInput) error {
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("input.Validate: %w", err)
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", domain.ErrUUIDInvalid)
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		profile, err := u.postgres.GetProfile(ctx, id)
		if err != nil {
			return fmt.Errorf("u.postgres.GetProfile: %w", err)
		}

		if profile.IsDeleted() {
			return fmt.Errorf("profile.IsDeleted: %w", domain.ErrNotFound)
		}

		newProfile, err := update(profile, input)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}

		if newProfile == profile {
			return domain.ErrNoChangesFound
		}

		err = u.postgres.UpdateProfile(ctx, newProfile)
		if err != nil {
			return fmt.Errorf("u.postgres.UpdateProfile: %w", err)
		}

		// Обновляем данные в Redis
		err = u.redis.SetCache(ctx, newProfile)
		if err != nil {
			log.Error().Err(err).Str("profileID", profile.ID.String()).Msg("cache: UpdateProfile: set cache")
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction.Wrap: %w", err)
	}

	return nil
}

func update(profile domain.Profile, input dto.UpdateProfileInput) (domain.Profile, error) {
	if input.Name != nil {
		profile.Name = domain.Name(*input.Name)
	}

	if input.Age != nil {
		profile.Age = domain.Age(*input.Age)
	}

	if input.Email != nil {
		profile.Contacts.Email = *input.Email
	}

	if input.Phone != nil {
		profile.Contacts.Phone = *input.Phone
	}

	err := profile.Validate()
	if err != nil {
		return profile, fmt.Errorf("profile.Validate: %w", err)
	}

	return profile, nil
}
