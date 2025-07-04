package usecase

import (
	"context"
	"fmt"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

// Внутренняя логика нашего приложения (все что происходит после получения input, до преобразования в output)

func (u *UseCase) CreateProfile(ctx context.Context, input dto.CreateProfileInput) (dto.CreateProfileOutput, error) {
	var output dto.CreateProfileOutput

	profile, err := domain.NewProfile(input.Name, input.Age, input.Email, input.Phone)
	if err != nil {
		return output, fmt.Errorf("domain.NewProfile: %w", err)
	}

	ctx, err = transaction.Begin(ctx)
	if err != nil {
		return output, fmt.Errorf("transaction.Begin: %w", err)
	}

	defer transaction.Rollback(ctx)

	err = u.postgres.CreateProfile(ctx, profile)
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreateProfile: %w", err)
	}

	property := domain.NewProperty(profile.ID, []string{"home", "primary"})

	err = u.postgres.CreateProperty(ctx, property)
	if err != nil {
		return output, fmt.Errorf("u.postgres.CreateProperty: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return output, fmt.Errorf("transaction.Commit: %w", err)
	}

	return dto.CreateProfileOutput{
		ID: profile.ID,
	}, nil
}
