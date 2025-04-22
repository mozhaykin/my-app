package usecase

import (
	"github.com/google/uuid"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

type Cache interface {
	Add(key uuid.UUID, p domain.Profile)
	Get(key uuid.UUID) (domain.Profile, error)
	Update(key uuid.UUID, p domain.Profile)
	Delete(key uuid.UUID)
	PrintAll()
}

type UseCase struct {
	cache Cache
}

func New(cache Cache) *UseCase {
	return &UseCase{
		cache: cache,
	}
}
