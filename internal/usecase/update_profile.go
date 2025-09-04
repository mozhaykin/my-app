package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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

	ctx, err = transaction.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction.Begin: %w", err)
	}

	defer transaction.Rollback(ctx)

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

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("transaction.Commit: %w", err)
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
