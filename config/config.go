package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/mozhaykin/my-app/internal/adapter/kafkaproducer"
	"github.com/mozhaykin/my-app/internal/controller/grpc"
	"github.com/mozhaykin/my-app/internal/controller/kafkaconsumer"
	"github.com/mozhaykin/my-app/internal/controller/worker"
	"github.com/mozhaykin/my-app/pkg/grpcclientv1"
	"github.com/mozhaykin/my-app/pkg/httpclientv1"
	"github.com/mozhaykin/my-app/pkg/httpclientv2"
	"github.com/mozhaykin/my-app/pkg/httpserver"
	"github.com/mozhaykin/my-app/pkg/logger"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/postgres"
	"github.com/mozhaykin/my-app/pkg/redisclient"
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
	OTEL              otel.Config
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
