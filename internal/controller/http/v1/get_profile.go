package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	input := dto.GetProfileInput{
		ID: chi.URLParam(r, "id"),
	}

	err := input.Validate()
	if err != nil {
		http.Error(w, "validate error", http.StatusBadRequest)

		return
	}

	output, err := h.usecase.GetProfile(input)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error().Err(err).Msg("v1 get_profile: h.usecase.GetProfile")

		return
	}

	render.JSON(w, output, http.StatusOK)
}
