package dto

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNameInvalid = errors.New("name is invalid")

type CreateProfileOutput struct {
	ID uuid.UUID `json:"id"`
}

type CreateProfileInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (i CreateProfileInput) Validate() error {
	if i.Name == "" {
		return ErrNameInvalid
	}

	return nil
}
