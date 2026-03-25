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
	"github.com/mozhaykin/my-app/internal/controller/kafkaconsumer"
	"github.com/mozhaykin/my-app/internal/controller/worker"
	"github.com/mozhaykin/my-app/pkg/httpclientv1"
	"github.com/mozhaykin/my-app/pkg/httpserver"
	"github.com/mozhaykin/my-app/pkg/logger"
	"github.com/mozhaykin/my-app/pkg/otel"
	"github.com/mozhaykin/my-app/pkg/postgres"
	"github.com/mozhaykin/my-app/pkg/redisclient"
)

// make up
// make test_integration_http_v1

var ctx = context.Background()

type CreateProfileRequest = httpclientv1.CreateProfileRequest

type GetProfilesRequest = httpclientv1.GetProfilesRequest

type UpdateProfileRequest = httpclientv1.UpdateProfileRequest

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	client      *httpclientv1.Client
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
		HTTPServer: httpserver.Config{
			Port: "8080",
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
		HTTPClientV1: httpclientv1.Config{
			Host: "localhost",
			Port: "8080",
		},
		KafkaProducer: kafkaproducer.Config{
			Addr: []string{"localhost:9094"},
		},
		KafkaConsumer: kafkaconsumer.Config{
			Addr:     []string{"localhost:9094"},
			Topic:    "mozhaykin-my-app-topic",
			Group:    "mozhaykin-my-app-group",
			Disabled: true, // Disable consumer in test!
		},
		OutboxKafkaWorker: worker.OutboxKafkaConfig{
			Limit: 10,
		},
	}

	logger.Init(c.Logger)
	otel.SilentModeInit()

	s.PrepareTestDB(c.Postgres)

	s.kafkaWriter = &kafka.Writer{
		Addr: kafka.TCP(c.KafkaProducer.Addr...),
	}

	go func() {
		err := app.Run(context.Background(), c)
		s.Require().NoError(err)
	}()

	s.client = httpclientv1.New(c.HTTPClientV1)

	time.Sleep(time.Second)
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
