package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/config"
	"github.com/mozhaykin/my-app/internal/adapter/kafkaproducer"
	"github.com/mozhaykin/my-app/internal/adapter/postgres"
	"github.com/mozhaykin/my-app/internal/adapter/redis"
	"github.com/mozhaykin/my-app/internal/controller/grpc"
	"github.com/mozhaykin/my-app/internal/controller/http"
	"github.com/mozhaykin/my-app/internal/controller/kafkaconsumer"
	"github.com/mozhaykin/my-app/internal/controller/worker"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/pkg/httpserver"
	"github.com/mozhaykin/my-app/pkg/metrics"
	pgpool "github.com/mozhaykin/my-app/pkg/postgres"
	"github.com/mozhaykin/my-app/pkg/redisclient"
	"github.com/mozhaykin/my-app/pkg/router"
	"github.com/mozhaykin/my-app/pkg/transaction"
)

func Run(ctx context.Context, c config.Config) error {
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("pgpool.New: %w", err)
	}

	transaction.Init(pgPool)

	redisClient, err := redisclient.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}

	entityMetrics := metrics.NewEntity()   // kafka consumer
	httpMetrics := metrics.NewHTTPServer() // основные метрики

	kafkaProducer := kafkaproducer.New(c.KafkaProducer, entityMetrics)

	uc := usecase.New(
		postgres.New(),
		redis.New(redisClient),
		kafkaProducer,
	)

	outboxKafkaWorker := worker.NewOutboxKafka(uc, c.OutboxKafkaWorker)

	kafkaConsumer := kafkaconsumer.New(c.KafkaConsumer, entityMetrics, uc)

	// Worker (просто пример воркера)
	someWorker, err := worker.NewSomeWorker(uc)
	if err != nil {
		return fmt.Errorf("worker.NewSomeWorker: %w", err)
	}

	// GRPC
	grpcServer, err := grpc.New(c.GRPCServer, uc)
	if err != nil {
		return fmt.Errorf("grpc.New: %w", err)
	}

	// HTTP
	r := router.New()
	http.ProfileRouter(r, uc, httpMetrics)
	httpServer := httpserver.New(r, c.HTTPServer)

	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Info().Msg("App got signal to stop")

	// Controllers close
	outboxKafkaWorker.Stop()
	kafkaConsumer.Close()
	someWorker.Stop()
	grpcServer.Close()
	httpServer.Close()

	// Adapters close
	redisClient.Close()
	kafkaProducer.Close()
	pgPool.Close()

	log.Info().Msg("App stopped!")

	return nil
}
