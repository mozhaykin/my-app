package v1

import (
	"encoding/json"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

// Получаем данные из запроса и передаем их в функцию обработчик в UseCase,
// получаем ответ (output) и передаем его пользователю

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
