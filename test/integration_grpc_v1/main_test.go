//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/kafkaproducer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/grpc"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/worker"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/grpcclientv1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
)

// Вначае поднимается база данных, потом запускаются тесты.
// make up
// make test_integration_grpc_v1

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	profile     *grpcclientv1.Client
	kafkaWriter *kafka.Writer
}

// Запускается один раз, до тестов (например для поднятия коннекшена к базе).
func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	s.ResetMigrations() // Удаление миграций и накатывание заново.

	// Config
	c := config.Config{
		App: config.App{
			Name:    "my-app",
			Version: "test",
		},
		GRPC: grpc.Config{
			Host: "localhost",
			Port: "50051",
		},
		Logger: logger.Config{
			AppName:       "my-app",
			AppVersion:    "test",
			Level:         "debug",
			PrettyConsole: true,
		},
		Postgres: postgres.Config{
			Host:     "localhost",
			Port:     "5432",
			User:     "login",
			Password: "pass",
			DBName:   "postgres",
		},
		KafkaProducer: kafkaproducer.Config{
			Addr: []string{"localhost:9094"},
		},
		OutboxKafka: worker.OutboxKafkaConfig{
			Limit: 10,
		},
	}

	logger.Init(c.Logger)

	// Kafka writer for direct produce messages
	s.kafkaWriter = &kafka.Writer{
		Addr: kafka.TCP(c.KafkaProducer.Addr...),
	}

	// Server
	go func() {
		err := app.Run(context.Background(), c)
		s.NoError(err)
	}()

	// Client
	client, err := grpcclientv1.New(grpcclientv1.Config{Host: "localhost", Port: "50051"})
	s.NoError(err)

	s.profile = client

	time.Sleep(time.Second) // Спим секунду, что горутина с сервером успела запуститься
}

// Не заданные функции. Иногда ими удобно пользоваться:
func (s *Suite) TearDownSuite() {} // Запускается один раз в вконце, после тестов (например закроет коннекшн к базе данных).

func (s *Suite) SetupTest() {} // Запускается перед каждым кейсом (например прогоняются миграции).

func (s *Suite) TearDownTest() {} // Запускается после каждого кейса (например очистить базу данных).
