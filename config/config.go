package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclient"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App      App
	HTTP     httpserver.Config
	Logger   logger.Config
	Postgres postgres.Config
	Client   httpclient.Config
}

func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env") // загружаем переменные окружения из файла .env
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	err = envconfig.Process("", &config) // записываем данные в конфиг из переменных окружения
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
