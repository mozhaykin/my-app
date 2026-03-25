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

// make up
// make test_integration_grpc_v1

var ctx = context.Background()

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

	var err error
	s.client, err = grpcclientv1.New(c.GRPSClientV1)
	s.Require().NoError(err)

	time.Sleep(time.Second)
}

func (s *Suite) SetupTest() {}

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

func (s *Suite) TearDownSuite() {
	if s.db != nil {
		_ = s.db.Close()
	}
}
