package v2

import (
	"context"
	"errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) DeleteProfileByID(ctx context.Context, request http_server.DeleteProfileByIDRequestObject) (
	http_server.DeleteProfileByIDResponseObject, error,
) {
	input := dto.DeleteProfileInput{
		ID: request.ID.String(),
	}

	err := h.usecase.DeleteProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return http_server.DeleteProfileByID404JSONResponse{Error: err.Error()}, nil

		default:
			return http_server.DeleteProfileByID400JSONResponse{Error: err.Error()}, nil
		}
	}

	return http_server.DeleteProfileByID204Response{}, nil
}
