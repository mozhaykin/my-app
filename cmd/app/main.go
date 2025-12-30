package main

import (
	"context"

	"github.com/rs/zerolog/log"

	_ "go.uber.org/automaxprocs"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel"
)

// В пакете main создаем новый конфиг, инициализируем сторонний логгер и запускаем основную функцию
// приложения app.Run, передавая в нее стандартный context.Background() и конфиг.

// Основная функция Run находится не в пакете main, для того, чтобы в интеграционных тестах можно было ее запустить.
// Импортировать пакет main в другие пакеты нельзя!

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	ctx := context.Background()

	err = otel.Init(ctx, c.OTEL)
	if err != nil {
		log.Fatal().Err(err).Msg("otel.Init")
	}

	defer otel.Close()

	err = app.Run(ctx, c)
	if err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
