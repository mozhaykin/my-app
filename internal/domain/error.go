package domain

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrAllFieldsAreEmpty = errors.New("all fields for update are empty")
	ErrUUIDInvalid       = errors.New("uuid is invalid")
	ErrNoChangesFound    = errors.New("no changes found")
)
