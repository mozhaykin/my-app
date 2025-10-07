package http

import (
	"github.com/go-chi/chi/v5"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/server"
	ver1 "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http/v1"
	ver2 "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http/v2"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
)

func ProfileRouter(r *chi.Mux, uc *usecase.UseCase) {
	v1 := ver1.New(uc)
	v2 := ver2.New(uc)

	r.Route("/amozhaykin/my-app/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/profile", v1.CreateProfile)
			r.Get("/profile/{id}", v1.GetProfile)
			r.Delete("/profile/{id}", v1.DeleteProfile)
			r.Put("/profile", v1.UpdateProfile)
		})

		r.Route("/v2", func(r chi.Router) {
			mux := server.NewStrictHandler(v2, []server.StrictMiddlewareFunc{})
			server.HandlerFromMux(mux, r)
		})
	})
}
