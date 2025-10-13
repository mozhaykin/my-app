package v2

import (
	"context"
	"errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/server"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
)

func (h *Handlers) GetProfileByID(ctx context.Context, request server.GetProfileByIDRequestObject) (
	server.GetProfileByIDResponseObject, error,
) {
	input := dto.GetProfileInput{
		ID: request.ID,
	}

	baggage.PutProfileID(ctx, input.ID)

	output, err := h.usecase.GetProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return server.GetProfileByID404JSONResponse{Error: err.Error()}, nil

		default:
			return server.GetProfileByID400JSONResponse{Error: err.Error()}, nil
		}
	}

	var profile server.GetProfileByID200JSONResponse

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
