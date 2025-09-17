package v2

import (
	"context"
	"errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http_server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) GetProfileByID(ctx context.Context, request http_server.GetProfileByIDRequestObject) (
	http_server.GetProfileByIDResponseObject, error,
) {
	input := dto.GetProfileInput{
		ID: request.ID.String(),
	}

	output, err := h.usecase.GetProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return http_server.GetProfileByID404JSONResponse{Error: err.Error()}, nil

		default:
			return http_server.GetProfileByID400JSONResponse{Error: err.Error()}, nil
		}
	}

	var profile http_server.GetProfileByID200JSONResponse

	profile.ID = output.ID
	profile.Name = string(output.Name)
	profile.Age = int(output.Age)
	profile.Contacts.Email = output.Contacts.Email
	profile.Contacts.Phone = output.Contacts.Phone
	profile.CreatedAt = output.CreatedAt
	profile.UpdatedAt = output.UpdatedAt
	profile.Status = int(output.Status)
	profile.Verified = output.Verified

	return profile, nil
}
