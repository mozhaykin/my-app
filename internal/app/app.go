package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/config"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/kafkaproducer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/adapter/postgres"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/grpc"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/http"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/kafkaconsumer"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/controller/worker"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/httpserver"
	pgpool "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/postgres"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/router"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/transaction"
)

func Run(ctx context.Context, c config.Config) error {
	// Создаем пул подключений к Postgres (используя данные из конфиг файла)
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("pgpool.New: %w", err)
	}

	transaction.Init(pgPool)

	// Kafka producer
	kafkaProducer := kafkaproducer.New(c.KafkaProducer)

	// UseCase
	uc := usecase.New(
		postgres.New(),
		kafkaProducer,
	)

	// Kafka consumer
	kafkaConsumer := kafkaconsumer.New(c.KafkaConsumer, uc)

	// Outbox Kafka worker
	outboxKafkaWorker := worker.NewOutboxKafka(uc, c.OutboxKafka)

	// GRPC
	grpcServer, err := grpc.New(c.GRPC, uc)
	if err != nil {
		return fmt.Errorf("grpc.New: %w", err)
	}

	// HTTP
	r := router.New()         // Создаем новый роутер chi
	http.ProfileRouter(r, uc) // Прописываем ручки
	// Создаем HTTP сервер, передавая в него роутер и используя данные из конфиг файла
	httpServer := httpserver.New(r, c.HTTP)

	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig // wait signal

	log.Info().Msg("App got signal to stop")

	// Controllers close
	grpcServer.Close()
	httpServer.Close()
	outboxKafkaWorker.Close()
	kafkaConsumer.Close()

	// Adapters close
	pgPool.Close()
	kafkaProducer.Close()

	log.Info().Msg("App stopped!")

	return nil
}
