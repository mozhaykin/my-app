//go:build integration

package test

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
)

func (s *Suite) PrepareTestDB(c postgres.Config) {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)

	filePath := "file://../../migration/postgres"

	db, err := sql.Open("postgres", dbURL)
	s.Require().NoError(err)

	// Полная очистка схемы перед миграцией и началом всех тестов
	_, err = db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	s.Require().NoError(err)

	// Сохраняем db в Suite, чтобы SetupTest или TearDownTest могли использовать её
	s.db = db

	// Накатываем миграции
	m, err := migrate.New(filePath, dbURL)
	s.Require().NoError(err)

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		err = nil
	}
	s.Require().NoError(err)
}
