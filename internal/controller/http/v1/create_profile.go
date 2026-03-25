package v1

import (
	"encoding/json"
	"net/http"

	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/dto/baggage"
	"github.com/mozhaykin/my-app/pkg/render"
)

func (h *Handlers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.CreateProfileInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "json decode error: ")

		return
	}

	output, err := h.usecase.CreateProfile(ctx, input)
	if err != nil {
		render.Error(ctx, w, err, http.StatusBadRequest, "request failed: ")

		return
	}

	baggage.PutProfileID(ctx, output.ID.String())

	render.JSON(w, output, http.StatusCreated)
}
