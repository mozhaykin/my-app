package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
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

	err = u.postgres.UpdateProfile(ctx, newProfile)
	if err != nil {
		return fmt.Errorf("u.postgres.UpdateProfile: %w", err)
	}

	return nil
}

func update(profile domain.Profile, input dto.UpdateProfileInput) (domain.Profile, error) {
	changes := HasChanges(profile, input)
	if !changes {
		return profile, domain.ErrNoChangesFound
	}

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

func HasChanges(p domain.Profile, input dto.UpdateProfileInput) bool {
	if input.Name != nil && p.Name != domain.Name(*input.Name) {
		return true
	}

	if input.Age != nil && p.Age != domain.Age(*input.Age) {
		return true
	}

	if input.Email != nil && p.Contacts.Email != *input.Email {
		return true
	}

	if input.Phone != nil && p.Contacts.Phone != *input.Phone {
		return true
	}

	return false
}
