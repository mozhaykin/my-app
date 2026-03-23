package dto

import "github.com/mozhaykin/my-app/internal/domain"

type GetProfileOutput struct {
	domain.Profile
}

type GetProfileInput struct {
	ID string
}
