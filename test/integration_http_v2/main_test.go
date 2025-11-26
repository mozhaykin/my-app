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
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/kafkaproducer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/kafkaconsumer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/worker"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv2"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
)

// make up 								поднимается база данных
// make test_integration_http_v2 		запускаются тесты

var ctx = context.Background()

type CreateProfileRequest = httpclientv2.CreateProfileRequest

type GetProfilesRequest = httpclientv2.GetProfilesRequest

type UpdateProfileRequest = httpclientv2.UpdateProfileRequest

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	client      *httpclientv2.Client
	kafkaWriter *kafka.Writer
	kafkaReader *kafka.Reader
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
		HTTPClientV2: httpclientv2.Config{
			Address: "http://localhost:8080/amozhaykin/my-app/api/v2",
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

	// Подключение к базе и миграции
	s.PrepareTestDB(c.Postgres)

	// Kafka writer
	s.kafkaWriter = &kafka.Writer{
		Addr: kafka.TCP(c.KafkaProducer.Addr...),
	}

	// Kafka reader
	s.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.KafkaConsumer.Addr,
		Topic:   c.KafkaConsumer.Topic,
		GroupID: c.KafkaConsumer.Group,
	})

	// Server
	go func() {
		err := app.Run(context.Background(), c)
		s.NoError(err)
	}()

	// Client V2
	var err error
	s.client, err = httpclientv2.New(c.HTTPClientV2)
	s.NoError(err)

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
