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

	property := domain.NewProperty(profile.ID, []string{"home", "primary"})

	kafkaMsg, err := profile.ToKafkaMsg("awesome-topic")
	if err != nil {
		return output, fmt.Errorf("profile.ToEvent: %w", err)
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

		// Дополнительная запись profile в таблицу Outbox (из которой читает воркер и гарантировано отправляет в Кафку)
		err = u.postgres.SaveOutboxKafka(ctx, kafkaMsg)
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
