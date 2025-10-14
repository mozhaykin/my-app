package v2

import (
	"context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
)

func (h *Handlers) CreateProfile(ctx context.Context, request server.CreateProfileRequestObject) (
	server.CreateProfileResponseObject, error,
) {
	input := dto.CreateProfileInput{
		Name:  request.Body.Name,
		Age:   request.Body.Age,
		Email: request.Body.Email,
		Phone: request.Body.Phone,
	}

	output, err := h.usecase.CreateProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		return server.CreateProfile400JSONResponse{Error: err.Error()}, nil
	}

	baggage.PutProfileID(ctx, output.ID.String())

	return server.CreateProfile201JSONResponse{
		ID: output.ID,
	}, nil
}
