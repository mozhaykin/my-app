package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) DeleteProfile(input dto.DeleteProfileInput) error {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	u.cache.Delete(id)

	return nil
}
