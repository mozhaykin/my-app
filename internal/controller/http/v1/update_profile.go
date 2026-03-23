package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/dto/baggage"
	"github.com/mozhaykin/my-app/pkg/render"
)

func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.UpdateProfileInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "json decode error")

		return
	}

	baggage.PutProfileID(ctx, input.ID)

	err = h.usecase.UpdateProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(ctx, w, err, http.StatusNotFound, "")

		default:
			render.Error(ctx, w, err, http.StatusBadRequest, "request failed: ")
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
