package v1

import "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"

type Handlers struct {
	usecase *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handlers {
	return &Handlers{
		usecase: uc,
	}
}
