package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (u *UseCase) GetProfile(input dto.GetProfileInput) (dto.GetProfileOutput, error) {
	var output dto.GetProfileOutput

	key, err := uuid.Parse(input.ID)
	if err != nil {
		return output, fmt.Errorf("uuid.Parse: %w", err)
	}

	p, err := u.cache.Get(key)
	if err != nil {
		return output, fmt.Errorf("u.cache.Get: %w", err)
	}

	return dto.GetProfileOutput{
		Name: p.Name,
		Age:  p.Age,
	}, nil
}
