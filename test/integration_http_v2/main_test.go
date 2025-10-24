//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpclientv2"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
)

// Вначае поднимается база данных, потом запускаются тесты.
// make up
// make test_integration_http_v2

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	profile *httpclientv2.Client
}

// Запускается один раз, до тестов (например для поднятия коннекшена к базе).
func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	s.ResetMigrations()

	// Config
	c := config.Config{
		App: config.App{
			Name:    "my-app",
			Version: "test",
		},
		HTTP: httpserver.Config{
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
	}

	logger.Init(c.Logger)

	// Server
	go func() {
		err := app.Run(context.Background(), c)
		s.NoError(err)
	}()

	// Client V2
	var err error
	s.profile, err = httpclientv2.New("http://localhost:8080/amozhaykin/my-app/api/v2")
	s.NoError(err)

	time.Sleep(time.Second)
}

// Не заданные функции. Оставил для наглядности:
func (s *Suite) TearDownSuite() {} // Запускается один раз в вконце, после тестов (например закроет коннекшн к базе данных).

func (s *Suite) SetupTest() {} // Запускается перед каждым кейсом (например прогоняются миграции).

func (s *Suite) TearDownTest() {} // Запускается после каждого кейса (например очистить базу данных).
