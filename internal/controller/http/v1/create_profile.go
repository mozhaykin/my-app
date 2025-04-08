package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	log.Info().Msg(fmt.Sprintf("input %+v", r.Body))

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "json error", http.StatusBadRequest)
		log.Error().Err(err).Msg("http v1 create_profile: json.NewDecoder.Decode")

		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, "validate error", http.StatusBadRequest)
		log.Error().Err(err).Msg("v1 create_profile: input.Validate")

		return
	}

	h.cache.Add(input.Name, input.Age)

	// UseCase

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
