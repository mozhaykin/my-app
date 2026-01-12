package kafkaconsumer

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/metrics"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/otel/tracer"
)

type Config struct {
	Addr     []string `envconfig:"KAFKA_CONSUMER_ADDR"     required:"true"`
	Topic    string   `default:"profiles.created"          envconfig:"KAFKA_CONSUMER_TOPIC"`
	Group    string   `default:"awesome-group"             envconfig:"KAFKA_CONSUMER_GROUP"`
	Disabled bool     `envconfig:"KAFKA_CONSUMER_DISABLED"`
}
type Consumer struct {
	config  Config
	reader  *kafka.Reader
	usecase *usecase.UseCase
	metrics *metrics.Entity
	stop    context.CancelFunc
	done    chan struct{}
}

func New(cfg Config, m *metrics.Entity, uc *usecase.UseCase) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          cfg.Addr,
		Topic:            cfg.Topic,
		GroupID:          cfg.Group,
		ErrorLogger:      logger.ErrorLogger(),
		ReadBatchTimeout: time.Second,
		CommitInterval:   100 * time.Millisecond,
	})

	ctx, stop := context.WithCancel(context.Background())

	c := &Consumer{
		config:  cfg,
		reader:  r,
		usecase: uc,
		metrics: m,
		stop:    stop,
		done:    make(chan struct{}),
	}

	if c.config.Disabled {
		log.Info().Msg("kafka consumer: disabled")

		return c
	}

	go c.run(ctx)

	return c
}

func (c *Consumer) run(ctx context.Context) {
	log.Info().Msg("kafka consumer: started")

	const consume = "consume"

	for {
		now := time.Now()

		// Читаем сообщение из kafka
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("kafka consumer: FetchMessage")

			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				break
			}

			continue
		}

		// Восстанавливаем контекст из Kafka headers
		msgCtx := otel.CtxFromKafkaHeaders(ctx, m.Headers)

		// Создаем span от извлеченного контекста
		msgCtx, span := tracer.Start(msgCtx, "kafka consume "+m.Topic,
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				semconv.MessagingSystemKafka,
				semconv.MessagingDestinationName(m.Topic),
				semconv.MessagingKafkaOffset(int(m.Offset)),
				semconv.MessagingConsumerGroupName(c.config.Group),
				semconv.MessagingKafkaMessageKey(string(m.Key)),
			),
		)

		// Обрабатываем сообщение
		err = c.usecase.Consume(msgCtx, m)
		if err != nil {
			c.metrics.Total(consume, metrics.Error) // Инкрементим счетчик с ошибками
			log.Error().Err(err).Msg("kafka consumer: some work failed")

			span.End()

			continue
		}

		// Коммитим оффсет в consumer group
		err = c.reader.CommitMessages(ctx, m)
		if err != nil {
			c.metrics.Total(consume, metrics.Error) // Инкрементим счетчик с ошибками
			log.Error().Err(err).Msg("kafka consumer: CommitMessages")

			span.End()

			continue
		}

		c.metrics.Duration(consume, now)     // Считаем и записываем в метрику продолжительность обработки запроса
		c.metrics.Total(consume, metrics.Ok) // Инкрементим счетчик со статусом OK.

		span.End() // Закрываем span
	}

	close(c.done)
}

func (c *Consumer) Close() {
	if c.config.Disabled {
		return
	}

	log.Info().Msg("kafka consumer: closing")

	c.stop()

	err := c.reader.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka consumer: reader.Close")
	}

	<-c.done

	log.Info().Msg("kafka consumer: closed")
}
