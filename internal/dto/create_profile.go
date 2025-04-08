package dto

import "errors"

var ErrNameInvalid = errors.New("name is invalid")

type CreateProfileInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (i *CreateProfileInput) Validate() error {
	if i.Name == "" {
		return ErrNameInvalid
	}

	return nil
}
