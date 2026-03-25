package usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

func (u *UseCase) SomeWork(ctx context.Context) error {
	_, span := tracer.Start(ctx, "usecase SomeWork")
	defer span.End()

	log.Info().Msg("SomeWork called")

	// Выполняется какая то работа, например вызова клиента другого сервиса.
	// Для этого в usecase.go создается интерфейс с методами. Интерфейс добавляется в структуру UseCase и его конструктор
	// Здесь показан синтетический пример вызова своего же клиента.
	//
	// p, err := u.profile.GetProfile(ctx, "8638341a-b68a-4291-84ee-94b147afeff9")
	// if err != nil {
	//	return fmt.Errorf("SomeWork: %w", err)
	//}
	//
	// fmt.Println(p)

	return nil
}
