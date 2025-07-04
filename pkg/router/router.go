package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	// Добавляем middleware Recoverer, который:
	// - ловит паники (panic) во время обработки запросов
	// - возвращает клиенту HTTP 500 вместо падения сервиса
	// - логирует ошибку
	r.Use(middleware.Recoverer)

	r.Get("/live", probe)  // для проверки "жив" ли сервис
	r.Get("/ready", probe) // для проверки готовности сервиса принимать трафик

	return r
}

func probe(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("path", r.URL.Path).Msg("probe") // Логируем факт обращения к эндпоинту (с указанием пути)

	w.WriteHeader(http.StatusNoContent)
}
