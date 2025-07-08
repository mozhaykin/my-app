package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	input := dto.GetProfileInput{
		ID: chi.URLParam(r, "id"), // Достаем ключ из запроса
	}

	output, err := h.usecase.GetProfile(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(w, fmt.Errorf("h.usecase.GetProfile: %w", err), http.StatusNotFound)

		default:
			render.Error(w, fmt.Errorf("h.usecase.GetProfile: %w", err), http.StatusBadRequest)
		}

		return
	}

	render.JSON(w, output, http.StatusOK)
}
