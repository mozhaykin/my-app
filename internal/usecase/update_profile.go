package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) UpdateProfile(input dto.UpdateProfileInput) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	profile, err := domain.NewProfile(input.Name, input.Age, id)
	if err != nil {
		return fmt.Errorf("domain.NewProfile: %w", err)
	}

	u.cache.Update(id, profile)

	return nil
}
