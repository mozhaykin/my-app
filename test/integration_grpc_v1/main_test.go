//go:build integration

package test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mozhaykin/my-app/config"
	"github.com/mozhaykin/my-app/internal/adapter/kafkaproducer"
	"github.com/mozhaykin/my-app/internal/app"
	"github.com/mozhaykin/my-app/internal/controller/grpc"
	"github.com/mozhaykin/my-app/internal/controller/kafkaconsumer"
	"github.com/mozhaykin/my-app/internal/controller/worker"
	"github.com/mozhaykin/my-app/pkg/grpcclientv1"
	"github.com/mozhaykin/my-app/pkg/logger"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/postgres"
	"github.com/mozhaykin/my-app/pkg/redisclient"
)

// make up 								поднимается база данных
// make test_integration_grpc_v1		запускаются тесты

var ctx = context.Background()

// Указываю типы здесь для удобства, чтобы все тесты и миграции были одинаковыми для всех протоколов,
// различия есть только в файлах main_test.go
type CreateProfileRequest = grpcclientv1.CreateProfileRequest

type UpdateProfileRequest = grpcclientv1.UpdateProfileRequest

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	client      *grpcclientv1.Client
	kafkaWriter *kafka.Writer
	db          *sql.DB
}

// Запускается один раз, до тестов (например для поднятия коннекшена к базе).
func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	// Config
	c := config.Config{
		App: config.App{
			Name:    "my-app",
			Version: "test",
		},
		GRPCServer: grpc.Config{
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
		Redis: redisclient.Config{
			Addr: "localhost:6379",
		},
		GRPSClientV1: grpcclientv1.Config{
			Host: "localhost",
			Port: "50051",
		},
		KafkaProducer: kafkaproducer.Config{
			Addr: []string{"localhost:9094"},
		},
		KafkaConsumer: kafkaconsumer.Config{
			Addr:     []string{"localhost:9094"},
			Topic:    "amozhaykin-my-app-topic",
			Group:    "amozhaykin-my-app-group",
			Disabled: true, // Disable consumer in test!
		},
		OutboxKafkaWorker: worker.OutboxKafkaConfig{
			Limit: 10,
		},
	}

	logger.Init(c.Logger)
	otel.SilentModeInit() // явно отключаем otel

	// Подключение к базе и миграции
	s.PrepareTestDB(c.Postgres)

	// Kafka writer
	s.kafkaWriter = &kafka.Writer{
		Addr: kafka.TCP(c.KafkaProducer.Addr...),
	}

	// Server
	go func() {
		err := app.Run(context.Background(), c)
		s.Require().NoError(err)
	}()

	// Client
	var err error
	s.client, err = grpcclientv1.New(c.GRPSClientV1)
	s.Require().NoError(err)

	time.Sleep(time.Second) // Спим секунду, что горутина с сервером успела запуститься
}

// Запускается перед каждым кейсом
func (s *Suite) SetupTest() {}

// Запускается после каждого кейса
func (s *Suite) TearDownTest() {
	// Очистка данных из всех таблиц. Автоматически обходит все таблицы и работает при любой структуре БД.
	_, err := s.db.Exec(`
		DO $$
		DECLARE
		    r RECORD;
		BEGIN
		    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
		        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
		    END LOOP;
		END $$;
	`)
	s.NoError(err)
}

// Запускается один раз в вконце, после тестов (например для закрытия коннекшн к базе данных)
func (s *Suite) TearDownSuite() {
	if s.db != nil {
		_ = s.db.Close()
	}
}
