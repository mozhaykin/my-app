package usecase

import (
	"context"
	"fmt"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func (u *UseCase) GetProfiles(ctx context.Context, input dto.GetProfilesInput) (dto.GetProfilesOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase GetProfiles")
	defer span.End()

	var output dto.GetProfilesOutput

	err := input.Validate()
	if err != nil {
		return output, fmt.Errorf("input.Validate: %w", err)
	}

	if input.Limit == 0 {
		input.Limit = 10
	}

	if input.Order == "" {
		input.Order = "asc"
	}

	// Будем считать что запрос на получение нескольких профилей сразу редкий, поэтому
	// в Redis данные не проверяем, а сразу идем в базу

	profiles, err := u.postgres.GetProfiles(ctx, input)
	if err != nil {
		return output, fmt.Errorf("u.postgres.GetProfiles: %w", err)
	}

	if len(profiles) == 0 {
		return output, domain.ErrNotFound
	}

	output.Profiles = profiles

	return output, nil
}
