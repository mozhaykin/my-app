package v2

import (
	"context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) CreateProfile(ctx context.Context, request http_server.CreateProfileRequestObject) (
	http_server.CreateProfileResponseObject, error,
) {
	input := dto.CreateProfileInput{
		Name:  request.Body.Name,
		Age:   request.Body.Age,
		Email: string(request.Body.Email),
		Phone: request.Body.Phone,
	}

	output, err := h.usecase.CreateProfile(ctx, input)
	if err != nil {
		return http_server.CreateProfile400JSONResponse{Error: err.Error()}, nil //nolint: nilerr
	}

	return http_server.CreateProfile201JSONResponse{
		ID: output.ID,
	}, nil
}
