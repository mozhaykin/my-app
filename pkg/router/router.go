package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/live", probe)
	r.Get("/ready", probe)

	return r
}

func probe(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("path", r.URL.Path).Msg("probe")

	w.WriteHeader(http.StatusNoContent)
}
