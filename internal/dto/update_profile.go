package dto

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

// Update это единственный кейс, когда простые типы передаются по указателю. Это делается для того,
// чтобы мы могли понять было передано поле или нет.

// Если бы поля передавались по значению, то в том случае, если пользователь передал не все поля, а только те,
// которые хочет обновить, то не переданные поля всеравно бы инициализировались пустыми значениями,
// которые потом не прошли бы валидацию и update закончился бы ошибкой, что неверно, т.к. пользователь
// просто хотел обновить не все поля.

// Таким образом, если пользователь ничего не передал, мы увидим значение поля nil и будем понимать,
// что это поле обновлять не нужно.

type UpdateProfileInput struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Age   *int    `json:"age"`
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

func (i UpdateProfileInput) Validate() error {
	if i.Name == nil && i.Age == nil && i.Email == nil && i.Phone == nil {
		return domain.ErrAllFieldsAreEmpty
	}

	return nil
}
