package v2

import (
	"context"
	"errors"

	"github.com/mozhaykin/my-app/gen/http/profile_v2/server"
	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/dto/baggage"
)

func (h *Handlers) UpdateProfile(ctx context.Context, request server.UpdateProfileRequestObject) (
	server.UpdateProfileResponseObject, error,
) {
	input := dto.UpdateProfileInput{
		ID:    request.Body.ID,
		Name:  request.Body.Name,
		Age:   request.Body.Age,
		Email: request.Body.Email,
		Phone: request.Body.Phone,
	}

	baggage.PutProfileID(ctx, input.ID)

	err := h.usecase.UpdateProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return server.UpdateProfile404JSONResponse{Error: err.Error()}, nil

		default:
			return server.UpdateProfile400JSONResponse{Error: err.Error()}, nil
		}
	}

	return server.UpdateProfile204Response{}, nil
}
