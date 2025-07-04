package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/postgres"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	pgpool "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/router"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"

	"github.com/rs/zerolog/log"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
)

func Run(ctx context.Context, c config.Config) (err error) {
	// Создаем пул подключений к Postgres (используя данные из конфиг файла)
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("pgpool.New: %w", err)
	}

	transaction.Init(pgPool)

	// Создаем структуру UseCase, которая содержит интерфейс с методами обращения к базе данных
	uc := usecase.New(postgres.New(pgPool))

	// HTTP
	r := router.New()                       // Создаем новый роутер chi
	http.ProfileRouter(r, uc)               // Прописываем ручки
	httpServer := httpserver.New(r, c.HTTP) // Создаем HTTP сервер, передавая в него роутер и используя данные из конфиг файла

	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig // wait signal

	log.Info().Msg("App got signal to stop")

	httpServer.Close() //nolint:contextcheck
	pgPool.Close()

	log.Info().Msg("App stopped!")

	return nil
}
