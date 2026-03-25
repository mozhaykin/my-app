package usecase

import (
	"context"
	"fmt"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func (u *UseCase) CreateProfile(ctx context.Context, input dto.CreateProfileInput) (dto.CreateProfileOutput, error) {
	ctx, span := tracer.Start(ctx, "usecase CreateProfile")
	defer span.End()

	var output dto.CreateProfileOutput

	profile, err := domain.NewProfile(input.Name, input.Age, input.Email, input.Phone)
	if err != nil {
		return output, fmt.Errorf("domain.NewProfile: %w", err)
	}

	property := domain.NewProperty(profile.ID, []string{"home", "primary"})

	event, err := domain.EventProfileCreated(profile)
	if err != nil {
		return output, fmt.Errorf("profile.ToProfileCreatedEvent: %w", err)
	}

	// Дополняю event трейсом
	event.TraceContext, err = otel.ExtractTraceContext(ctx)
	if err != nil {
		return output, fmt.Errorf("extract trace context: %w", err)
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		err = u.postgres.CreateProfile(ctx, profile)
		if err != nil {
			return fmt.Errorf("u.postgres.CreateProfile: %w", err)
		}

		err = u.postgres.CreateProperty(ctx, property)
		if err != nil {
			return fmt.Errorf("u.postgres.CreateProperty: %w", err)
		}

		err = u.postgres.SaveOutbox(ctx, event)
		if err != nil {
			return fmt.Errorf("u.postgres.SaveOutboxKafka: %w", err)
		}

		return nil
	})
	if err != nil {
		return output, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return dto.CreateProfileOutput{
		ID: profile.ID,
	}, nil
}
