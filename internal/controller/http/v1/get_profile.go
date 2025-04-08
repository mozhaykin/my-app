package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog/log"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	input := chi.URLParam(r, "username")

	output, ok := h.cache.Get(input)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		log.Info().Msg("http v1 get_profile: h.cache.Get: key not found")

		return
	}

	data, err := json.Marshal(&output)
	if err != nil {
		http.Error(w, "json error", http.StatusBadRequest)
		log.Error().Err(err).Msg("http v1 get_profile: json.Marshal")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "write error", http.StatusBadRequest)
		log.Error().Err(err).Msg("http v1 get_profile: w.Write")

		return
	}
}
