package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/kafkaproducer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/grpc"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/kafkaconsumer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/worker"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/grpcclientv1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv2"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/redisclient"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App               App
	HTTPServer        httpserver.Config
	GRPCServer        grpc.Config
	Logger            logger.Config
	Postgres          postgres.Config
	Redis             redisclient.Config
	GRPSClientV1      grpcclientv1.Config
	HTTPClientV1      httpclientv1.Config
	HTTPClientV2      httpclientv2.Config
	KafkaProducer     kafkaproducer.Config
	KafkaConsumer     kafkaconsumer.Config
	OutboxKafkaWorker worker.OutboxKafkaConfig
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
