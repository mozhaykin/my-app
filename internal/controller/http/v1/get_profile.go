package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.GetProfileInput{
		ID: chi.URLParam(r, "id"), // Достаем ключ из запроса
	}

	baggage.PutProfileID(ctx, input.ID)

	output, err := h.usecase.GetProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(ctx, w, err, http.StatusNotFound, "")

		default:
			render.Error(ctx, w, err, http.StatusBadRequest, "request failed: ")
		}

		return
	}

	render.JSON(w, output, http.StatusOK)
}
