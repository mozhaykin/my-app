package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

type OutboxKafkaConfig struct {
	Limit int `default:"10" envconfig:"OUTBOX_KAFKA_WORKER_LIMIT"` // максимум 10 сообщений за один раз
}

type OutboxKafkaWorker struct {
	config  OutboxKafkaConfig
	usecase *usecase.UseCase
	stop    chan struct{} // канал для команды стоп
	done    chan struct{} // канал для получения сигнала о том, что стоп выполнен
}

func NewOutboxKafka(uc *usecase.UseCase, c OutboxKafkaConfig) *OutboxKafkaWorker {
	w := &OutboxKafkaWorker{
		config:  c,
		usecase: uc,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}

	go w.run()

	return w
}

func (w *OutboxKafkaWorker) run() {
	log.Info().Msg("outbox kafka worker: started")

FOR:
	for {
		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "worker outbox kafka", trace.WithSpanKind(trace.SpanKindInternal))

		eventsCount, err := w.usecase.OutboxReadAndProduce(ctx, w.config.Limit)
		if err != nil {
			log.Error().Err(err).Msg("outbox kafka worker: read and produce failed")
		}

		log.Info().Int("count", eventsCount).Msg("outbox kafka worker: read and produce")

		span.End()

		var sleepDuration time.Duration

		if eventsCount < w.config.Limit {
			sleepDuration = 10 * time.Second

			log.Info().Msg("outbox kafka worker: sleeping 10s")
		}

		select {
		case <-w.stop:
			break FOR
		case <-time.After(sleepDuration):
		}
	}

	log.Info().Msg("outbox kafka worker: stopped")

	close(w.done)
}

func (w *OutboxKafkaWorker) Stop() {
	log.Info().Msg("outbox kafka worker: closing")

	close(w.stop)

	<-w.done

	log.Info().Msg("outbox kafka worker: closed")
}
