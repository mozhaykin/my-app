//go:build integration

package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/kafkaproducer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/kafkaconsumer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/worker"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
)

// make up 								поднимается база данных
// migrate-up 							накатываются миграции
// make test_integration_http_v1 		запускаются тесты

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	profile     *httpclientv1.Client
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
		HTTPClientV1: httpclientv1.Config{
			Host: "localhost",
			Port: "8080",
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

	// Подключение к базе для очистки таблиц перед каждым тестом
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	s.NoError(err)
	s.db = db

	err = s.db.Ping()
	s.NoError(err)

	// Kafka writer
	s.kafkaWriter = &kafka.Writer{
		Addr: kafka.TCP(c.KafkaProducer.Addr...),
	}

	// Server с приложением
	go func() {
		err := app.Run(context.Background(), c)
		s.NoError(err)
	}()

	// Client V1
	s.profile = httpclientv1.New(c.HTTPClientV1)

	time.Sleep(time.Second) // Спим секунду, что горутина с сервером успела запуститься
}

// Запускается перед каждым кейсом
func (s *Suite) SetupTest() {
	// Очистка всех таблиц. Автоматически обходит все таблицы схемы и работает при любой структуре БД.
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

// Запускается после каждого кейса (например очистить базу данных)
func (s *Suite) TearDownTest() {}

// Запускается один раз в вконце, после тестов (например закроет коннекшн к базе данных)
func (s *Suite) TearDownSuite() {
	if s.db != nil {
		_ = s.db.Close()
	}
}
