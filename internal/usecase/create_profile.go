package usecase

import (
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) CreateProfile(input dto.CreateProfileInput) (dto.CreateProfileOutput, error) {
	var output dto.CreateProfileOutput

	profile, err := domain.NewProfile(input.Name, input.Age)
	if err != nil {
		return output, fmt.Errorf("domain.NewProfile: %w", err)
	}

	u.cache.Add(profile.ID, profile)

	return dto.CreateProfileOutput{
		ID: profile.ID,
	}, nil
}
