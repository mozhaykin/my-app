package domain

import (
	"errors"
)

var (
	ErrNotFound                = errors.New("not found")
	ErrAllFieldsForUpdateEmpty = errors.New("all fields for update are empty")
	ErrUUIDInvalid             = errors.New("uuid is invalid")
	ErrUUIDIsEmpty             = errors.New("uuid is empty")
	ErrNoChangesFound          = errors.New("no changes found")
)
