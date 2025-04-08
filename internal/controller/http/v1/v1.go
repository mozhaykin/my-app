package v1

import "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"

type Handlers struct {
	cache *dto.Cache
}

func New() *Handlers {
	return &Handlers{
		cache: dto.NewCache(),
	}
}
