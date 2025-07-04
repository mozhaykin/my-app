package dto

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

type UpdateProfileInput struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

func (i UpdateProfileInput) Validate() error {
	if i.ID == "" {
		return domain.ErrUUIDIsEmpty
	}

	if *i.Name == "" && *i.Age == 0 && *i.Email == "" && *i.Phone == "" {
		return domain.ErrAllFieldsForUpdateEmpty
	}

	return nil
}
