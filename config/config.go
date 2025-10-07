package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/grpc"
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
	GRPC     grpc.Config
	Logger   logger.Config
	Postgres postgres.Config
}

func New() (Config, error) {
	var config Config

	// При локальном запуске: загружаем переменные окружения рантайма из файла .env
	// При запуске на проде: считываем переменные окружения из среды (хранятся в секретах гитлаба)
	// и записываем в файл .env, затем загружаем переменные окружения рантайма из файла .env
	// Load не перезаписывает уже существующие переменные окружения, а только добавляет новые.
	err := godotenv.Load(".env")
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	// Записываем данные в структуру config из переменных окружения рантайма
	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
