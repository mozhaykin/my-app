package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	input := dto.UpdateProfileInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, fmt.Errorf("json.NewDecoder.Decode: %w", err), http.StatusBadRequest)

		return
	}

	err = h.usecase.UpdateProfile(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(w, fmt.Errorf("h.usecase.UpdateProfile: %w", err), http.StatusNotFound)

		default:
			render.Error(w, fmt.Errorf("h.usecase.UpdateProfile: %w", err), http.StatusBadRequest)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
