package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	input := dto.GetProfileInput{
		ID: chi.URLParam(r, "id"),
	}

	output, err := h.usecase.GetProfile(input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(w, err, http.StatusNotFound, "not found")

		default:
			render.Error(w, err, http.StatusBadRequest, "request failed")
		}

		return
	}

	render.JSON(w, output, http.StatusOK)
}
