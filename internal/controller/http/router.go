package http

import (
	"github.com/go-chi/chi/v5"
	ver1 "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http/v1"
)

func ProfileRouter(r *chi.Mux) {
	r.Route("/amozhaykin/my-app/api", func(r chi.Router) {
		v1 := ver1.New()

		r.Route("/v1", func(r chi.Router) {
			r.Get("/live", v1.Probe)
			r.Get("/ready", v1.Probe)
			r.Get("/profile/{username}", v1.GetProfile)
			r.Post("/profile", v1.CreateProfile)
		})
	})
}
