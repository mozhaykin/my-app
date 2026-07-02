package v2

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mozhaykin/my-app/gen/http/profile_v2/server"
	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
)

func (h *Handlers) GetProfiles(ctx context.Context, request server.GetProfilesRequestObject,
) (server.GetProfilesResponseObject, error) {
	input := dto.GetProfilesInput{
		Sort: request.Params.Sort,
	}

	if request.Params.Order != nil {
		input.Order = *request.Params.Order
	}

	if request.Params.Offset != nil {
		input.Offset = *request.Params.Offset
	}

	if request.Params.Limit != nil {
		input.Limit = *request.Params.Limit
	}

	output, err := h.usecase.GetProfiles(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return server.GetProfiles404JSONResponse{Error: err.Error()}, nil

		default:
			return server.GetProfiles400JSONResponse{Error: err.Error()}, nil
		}
	}

	profiles := make(server.GetProfiles200JSONResponse, 0, len(output.Profiles))

	for _, profile := range output.Profiles {
		var p server.GetProfileOutput

		p.ID = profile.ID
		p.CreatedAt = profile.CreatedAt
		p.UpdatedAt = profile.UpdatedAt
		p.Name = string(profile.Name)
		p.Age = int(profile.Age)
		p.Status = int(profile.Status)
		p.Verified = profile.Verified
		p.Contacts.Email = profile.Contacts.Email
		p.Contacts.Phone = profile.Contacts.Phone

		profiles = append(profiles, p)
	}

	return profiles, nil
}
