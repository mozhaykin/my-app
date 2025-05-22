export # Эта строка экспортирует все переменные, определенные в Makefile, в окружение, чтобы они были доступны в командах, которые выполняются внутри Makefile

DB_MIGRATE_URL = postgres://login:pass@localhost:5432/amozhaykin?sslmode=disable
MIGRATE_PATH = ./migration/postgres

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run: mod
	CGO_ENABLED=0 go run ./cmd/app

.PHONY: mod
mod:
	go mod tidy
	go mod download

.PHONY: mod-update
mod-update:
	go get -u all
	go mod tidy
	go mod download

.PHONY: up
up:
	docker compose up --build --force-recreate

.PHONY: down
down:
	docker compose down

.PHONY: down-v # Если хотим удалить контейнер вместе с volumes (удалится база данных)
down-v:
	docker compose down -v

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: integration-test
integration-test:
	go test -count=1 -v -tags integration ./test/integration

.PHONY: migrate-install
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

.PHONY: migrate-create
migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

.PHONY: migrate-up
migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

.PHONY: migrate-force-up # если dirty = true
migrate-force-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" force 20250519093401 # вставить свою версию схемы миграции

.PHONY: migrate-down
migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

.PHONY: seeder
seeder:
	go run ./cmd/seeder

#.PHONY: goose-status
#goose-status:
#	goose -dir=migrations postgres $(DATABASE_URL) status
#
#.PHONY: goose-up-one
#goose-up-one:
#	goose -dir=migrations postgres $(DATABASE_URL) up-by-one
#
#.PHONY: goose-down-one
#goose-down-one:
#	goose -dir=migrations postgres $(DATABASE_URL) down
#
#.PHONY: goose-up-all
#goose-up-all:
#	goose -dir=migrations postgres $(DATABASE_URL) up
#
#.PHONY: goose-reset
#goose-reset:
#	goose -dir=migrations postgres $(DATABASE_URL) reset