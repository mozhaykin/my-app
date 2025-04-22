package app

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/cache"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/router"

	"github.com/rs/zerolog/log"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
)

func Run() (err error) {
	uc := usecase.New(
		cache.New(),
	)

	r := router.New()

	http.ProfileRouter(r, uc)

	httpServer := httpserver.New(r, "8080")
	defer httpServer.Close()

	waiting(httpServer)

	return nil
}

func waiting(httpServer *httpserver.Server) {
	log.Info().Msg("App started!")

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)

	select {
	case i := <-wait:
		log.Info().Msg("App got signal: " + i.String())
	case err := <-httpServer.Notify():
		log.Error().Err(err).Msg("App got notify: httpServer.Notify")
	}

	log.Info().Msg("App is stopping...")
}
