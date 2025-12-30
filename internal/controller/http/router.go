package http //nolint:revive

import (
	"github.com/go-chi/chi/v5"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/http/profile_v2/server"
	ver1 "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http/v1"
	ver2 "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http/v2"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/metrics"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel"
)

func ProfileRouter(r *chi.Mux, uc *usecase.UseCase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)
	v2 := ver2.New(uc)

	r.Route("/amozhaykin/my-app/api", func(r chi.Router) {
		r.Use(logger.Middleware)
		r.Use(metrics.NewMiddleware(m))
		r.Use(otel.Middleware)

		r.Route("/v1", func(r chi.Router) {
			r.Post("/profile", v1.CreateProfile)
			r.Get("/profile/{id}", v1.GetProfile)
			r.Get("/profiles", v1.GetProfiles) // Ручка с пагинацией
			r.Put("/profile", v1.UpdateProfile)
			r.Delete("/profile/{id}", v1.DeleteProfile)
		})

		r.Route("/v2", func(r chi.Router) {
			mux := server.NewStrictHandler(v2, []server.StrictMiddlewareFunc{})
			server.HandlerFromMux(mux, r)
		})
	})
}
