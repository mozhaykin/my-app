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
	// Создаем пул подключений к Postgres (используя данные из конфиг файла)
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("pgpool.New: %w", err)
	}

	transaction.Init(pgPool)

	// Redis
	redisClient, err := redisclient.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}

	// Инициализация сервера с метриками.
	entityMetrics := metrics.NewEntity()   // Для kafka consumer
	httpMetrics := metrics.NewHTTPServer() // Основные метрики

	kafkaProducer := kafkaproducer.New(c.KafkaProducer, entityMetrics)

	// Создаем UseCase (передаем в структуру интерфейсы с методами которые вызываются в юзкейсах)
	uc := usecase.New(
		postgres.New(),
		redis.New(redisClient),
		kafkaProducer,
	)

	// Запускаем Outbox Kafka worker, в котором вызывается метод usecase OutboxReadAndProduce записывающий в kafka
	outboxKafkaWorker := worker.NewOutboxKafka(uc, c.OutboxKafkaWorker)

	// Запускаем kafka consumer, который читает сообщения из kafka, обрабатывает его вызывая метод usecase, и делает commit
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
	r := router.New()                      // Создаем новый роутер chi
	http.ProfileRouter(r, uc, httpMetrics) // Прописываем ручки
	// Создаем HTTP сервер, передавая в него роутер и используя данные из конфиг файла
	httpServer := httpserver.New(r, c.HTTPServer)

	log.Info().Msg("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig // wait signal

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
