package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string `envconfig:"HTTP_PORT" default:"8080"`
}

type Server struct { // Обёртка над стандартным http.Server для добавления методов start() и Close()
	server *http.Server
}

func New(handler http.Handler, c Config) *Server { // Создает экземпляр http.Server с настройками таймаутов
	httpServer := &http.Server{
		Handler:      handler,          // Обработчик запросов (например, роутер chi/gin)
		ReadTimeout:  20 * time.Second, // Таймаут на чтение запроса
		WriteTimeout: 20 * time.Second, // Таймаут на запись ответа
		Addr:         net.JoinHostPort("", c.Port),
	}

	s := &Server{
		server: httpServer,
	}

	go s.start() // Запуск сервера в горутине

	log.Info().Msg("http server: started on port: " + c.Port)

	return s
}

func (s *Server) start() {
	// Запускает сервер в блокирующем режиме (ListenAndServe).
	// Игнорирует ошибку http.ErrServerClosed (возникает при graceful shutdown).
	// Логирует другие ошибки (например, если порт занят).
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http server: ListenAndServe")
	}
}

func (s *Server) Close() {
	// Пытается остановить сервер, давая 25 секунд на завершение активных соединений.
	// Логирует ошибку, если shutdown не удался.
	// Всегда пишет лог о закрытии сервера.
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("http server: s.server.Shutdown")
	}

	log.Info().Msg("http server: closed")
}
