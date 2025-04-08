package v1

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	var input string

	in := r.URL.RawQuery

	err := json.Unmarshal([]byte(in), &input)
	if err != nil {
		http.Error(w, "json error", http.StatusBadRequest)
		log.Error().Err(err).Msg("http v1 get_profile: json.NewDecoder.Decode")

		return
	}

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
