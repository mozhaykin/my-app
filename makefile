export # Эта строка экспортирует все переменные, определенные в Makefile, в окружение, чтобы они были доступны в командах, которые выполняются внутри Makefile

#DATABASE_URL = postgres://user:pass@localhost:5432/postgres?sslmode=disable

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

#.PHONY: test
#test:
#	go test -v -cover -race ./...
#
#.PHONY: compose-up
#compose-up:
#	docker-compose -f "docker-compose.dev.yml" up --build -d --force-recreate
#	docker-compose -f "docker-compose.dev.yml" logs -f
#
#.PHONY: compose-down
#compose-down:
#	docker-compose -f "docker-compose.dev.yml" down --remove-orphans
#
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
#
#.PHONY: migrate-up
#migrate-up:
#	migrate -source file://new_migrations -database $(DATABASE_URL) up
#
#.PHONY: migrate-down
#migrate-down:
#	migrate -source file://new_migrations -database $(DATABASE_URL) down
#
#.PHONY: seeder
#seeder:
#	go run ./cmd/seeder