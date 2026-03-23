package http //nolint:revive

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mozhaykin/my-app/gen/http/profile_v2/server"
	ver1 "github.com/mozhaykin/my-app/internal/controller/http/v1"
	ver2 "github.com/mozhaykin/my-app/internal/controller/http/v2"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/pkg/logger"
	"github.com/mozhaykin/my-app/pkg/metrics"
	"github.com/mozhaykin/my-app/pkg/otel"
)

func ProfileRouter(r *chi.Mux, uc *usecase.UseCase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)
	v2 := ver2.New(uc)

	r.Handle("/metrics", promhttp.Handler())

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
