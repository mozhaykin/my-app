package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"

	"github.com/rs/zerolog/log"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

// Обязанности контроллера:
// Принять запрос в удобном для пользователя виде
// Сделать валидацию
// Передать дальше, получить ответ + обработать ошибку если есть
// Отдать в удобном для пользователя виде

func (h *Handlers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateProfileInput{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "json decode error", http.StatusBadRequest)
		log.Error().Err(err).Msg("http v1 create_profile: json.NewDecoder.Decode")

		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, "validate error", http.StatusBadRequest)
		log.Error().Err(err).Msg("v1 create_profile: input.Validate")

		return
	}

	output, err := h.usecase.CreateProfile(input)
	if err != nil {
		if errors.Is(err, domain.ErrAgeLessThan18) {
			err = errors.Unwrap(err)

			http.Error(w, "validate error: "+err.Error(), http.StatusBadRequest)
			log.Error().Err(err).Msg("v1 create_profile: h.usecase.CreateProfile")

			return
		}

		http.Error(w, "bad request", http.StatusBadRequest)
		log.Error().Err(err).Msg("v1 create_profile: h.usecase.CreateProfile")

		return
	}

	render.JSON(w, output, http.StatusOK)
}
