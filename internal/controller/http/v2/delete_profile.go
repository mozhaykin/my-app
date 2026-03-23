package v2

import (
	"context"
	"errors"

	"github.com/mozhaykin/my-app/gen/http/profile_v2/server"
	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/dto/baggage"
)

func (h *Handlers) DeleteProfileByID(ctx context.Context, request server.DeleteProfileByIDRequestObject) (
	server.DeleteProfileByIDResponseObject, error,
) {
	input := dto.DeleteProfileInput{
		ID: request.ID,
	}

	baggage.PutProfileID(ctx, input.ID)

	err := h.usecase.DeleteProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return server.DeleteProfileByID404JSONResponse{Error: err.Error()}, nil

		default:
			return server.DeleteProfileByID400JSONResponse{Error: err.Error()}, nil
		}
	}

	return server.DeleteProfileByID204Response{}, nil
}
