package v1

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (h *Handlers) Probe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	log.Info(). // логирование в хендлерах
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("Probe handler called")
}
