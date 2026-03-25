package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	// Добавляем middleware Recoverer, который:
	// - ловит паники (panic) во время обработки запросов
	// - возвращает клиенту HTTP 500 вместо падения сервиса
	// - логирует ошибку
	r.Use(middleware.Recoverer)

	r.Get("/live", probe)
	r.Get("/ready", probe)

	return r
}

func probe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
