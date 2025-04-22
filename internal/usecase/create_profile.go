package usecase

import (
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"

	"github.com/google/uuid"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) CreateProfile(input dto.CreateProfileInput) (dto.CreateProfileOutput, error) {
	var output dto.CreateProfileOutput

	key, err := uuid.NewUUID()
	if err != nil {
		return output, fmt.Errorf("uuid.NewUUID: %w", err)
	}

	profile, err := domain.NewProfile(input.Name, input.Age)
	if err != nil {
		return output, fmt.Errorf("domain.NewProfile: %w", err)
	}

	u.cache.Add(key, profile)

	return dto.CreateProfileOutput{
		ID: key,
	}, nil
}
