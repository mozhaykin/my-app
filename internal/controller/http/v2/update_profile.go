package v2

import (
	"context"
	"errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) UpdateProfile(ctx context.Context, request server.UpdateProfileRequestObject) (
	server.UpdateProfileResponseObject, error,
) {
	input := dto.UpdateProfileInput{
		ID:    request.Body.ID.String(),
		Name:  request.Body.Name,
		Age:   request.Body.Age,
		Email: request.Body.Email,
		Phone: request.Body.Phone,
	}

	err := h.usecase.UpdateProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return server.UpdateProfile404JSONResponse{Error: err.Error()}, nil

		default:
			return server.UpdateProfile400JSONResponse{Error: err.Error()}, nil
		}
	}

	return server.UpdateProfile204Response{}, nil
}
