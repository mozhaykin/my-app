package domain

import (
	"errors"
)

var (
	ErrEmptyName     = errors.New("name is empty")
	ErrEmptyID       = errors.New("id is empty")
	ErrAgeLessThan18 = errors.New("age is less than 18")
	ErrNotFound      = errors.New("not found")
)
