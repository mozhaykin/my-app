package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

// Получаем данные из запроса и передаем их в функцию обработчик в UseCase,
// получаем ответ (output) и передаем его пользователю

func (h *Handlers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateProfileInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		render.Error(w, fmt.Errorf("json.NewDecoder.Decode: %w", err), http.StatusBadRequest)

		return
	}

	output, err := h.usecase.CreateProfile(r.Context(), input)
	if err != nil {
		render.Error(w, fmt.Errorf("h.usecase.CreateProfile: %w", err), http.StatusBadRequest)

		return
	}

	render.JSON(w, output, http.StatusCreated)
}
