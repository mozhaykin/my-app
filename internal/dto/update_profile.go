package dto

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

// Практически единственный кейс, когда простые типы передаются по указателю. Это делается для того,
// чтобы мы могли различить следующие варианты:
// пользователь передал пустое значение, на которое нужно обновить.
// пользователь ничего не передал (nil) в этом случае поле обновлять не надо.
// Если бы поля в структуре были указаны по значению, то мы бы не смогли отличить нулевое значение от nil.
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

	if i.Name == nil && i.Age == nil && i.Email == nil && i.Phone == nil {
		return domain.ErrAllFieldsForUpdateEmpty
	}

	return nil
}
