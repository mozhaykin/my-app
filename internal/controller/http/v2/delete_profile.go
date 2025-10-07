package v2

import (
	"context"
	"errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) DeleteProfileByID(ctx context.Context, request server.DeleteProfileByIDRequestObject) (
	server.DeleteProfileByIDResponseObject, error,
) {
	input := dto.DeleteProfileInput{
		ID: request.ID.String(),
	}

	err := h.usecase.DeleteProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return server.DeleteProfileByID404JSONResponse{Error: err.Error()}, nil

		default:
			return server.DeleteProfileByID400JSONResponse{Error: err.Error()}, nil
		}
	}

	return server.DeleteProfileByID204Response{}, nil
}
