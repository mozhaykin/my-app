package dto

import "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"

type GetProfileOutput struct {
	Name domain.Name `json:"name"`
	Age  domain.Age  `json:"age"`
}

type GetProfileInput struct {
	ID string
}

func (i GetProfileInput) Validate() error {
	if i.ID == "" {
		return domain.ErrEmptyID
	}

	return nil
}
