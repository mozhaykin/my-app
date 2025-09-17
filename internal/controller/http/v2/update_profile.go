package v2

import (
	"context"
	"errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) UpdateProfile(ctx context.Context, request http_server.UpdateProfileRequestObject) (
	http_server.UpdateProfileResponseObject, error,
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
			return http_server.UpdateProfile404JSONResponse{Error: err.Error()}, nil

		default:
			return http_server.UpdateProfile400JSONResponse{Error: err.Error()}, nil
		}
	}

	return http_server.UpdateProfile204Response{}, nil
}
