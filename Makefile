export # Эта строка экспортирует все переменные, определенные в Makefile, в окружение, чтобы они были доступны в командах, которые выполняются внутри Makefile

DB_MIGRATE_URL = postgres://login:pass@localhost:5432/postgres?sslmode=disable
MIGRATE_PATH = ./migration/postgres

run: mod
	CGO_ENABLED=0 go run ./cmd/app

lint:
	golangci-lint run

mod:
	go mod tidy
	go mod download

mod-update:
	go get -u all
	go mod tidy
	go mod download

seeder:
	go run ./cmd/seeder

up:
	docker compose up --build --force-recreate

down:
	docker compose down

down-v: # Если хотим удалить контейнер вместе с volumes (удалится база данных)
	docker compose down -v

.PHONY: test
test:
	go test -v -cover ./...

integration_test_http_v1:
	go test -count=1 -v -tags integration ./test/integration_http_v1

integration_test_http_v2:
	go test -count=1 -v -tags integration ./test/integration_http_v2

integration_test_grpc_v1:
	go test -count=1 -v -tags integration ./test/integration_grpc_v1

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-force-up: # если dirty = true
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" force 20250519093401 # вставить свою версию схемы миграции

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

oapi-install:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

generate:
	go generate ./...

grpc_gen:
	mkdir -p ./gen/grpc/profile_v1
	protoc \
	  --proto_path=api/grpc \
	  --go_out=./gen/grpc/profile_v1 --go_opt=paths=source_relative \
	  --go-grpc_out=./gen/grpc/profile_v1 --go-grpc_opt=paths=source_relative \
	  profile_v1.proto

#goose-status:
#	goose -dir=migrations postgres $(DATABASE_URL) status
#
#goose-up-one:
#	goose -dir=migrations postgres $(DATABASE_URL) up-by-one
#
#goose-down-one:
#	goose -dir=migrations postgres $(DATABASE_URL) down
#
#goose-up-all:
#	goose -dir=migrations postgres $(DATABASE_URL) up
#
#goose-reset:
#	goose -dir=migrations postgres $(DATABASE_URL) reset