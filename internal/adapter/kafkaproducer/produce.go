package kafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/pkg/logger"
	"github.com/mozhaykin/my-app/pkg/metrics"
	"github.com/mozhaykin/my-app/pkg/otel/tracer"
)

type Config struct {
	Addr  []string `envconfig:"KAFKA_WRITER_ADDR"  required:"true"`
	Topic string   `envconfig:"KAFKA_WRITER_TOPIC"`
}

type Producer struct {
	config  Config
	writer  *kafka.Writer
	metrics *metrics.Entity
}

func New(c Config, m *metrics.Entity) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.Addr...),
		Balancer:     &kafka.Hash{Hasher: fnv.New32a()},
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  logger.ErrorLogger(),
		Async:        false, // можно включить отправку батчами (с интервалом в 1s) если producer захлёбывается
		// BatchSize:    10,
		// BatchTimeout: 100 * time.Millisecond,
	}

	return &Producer{
		config:  c,
		writer:  w,
		metrics: m,
	}
}

func (p *Producer) Produce(ctx context.Context, events []domain.Event) error {
	ctx, span := tracer.Start(ctx, "adapter kafka Produce", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	const produce = "produce"

	eventsCount := len(events)

	defer p.metrics.Duration(produce, time.Now())

	messages, err := buildKafkaMessages(events)
	if err != nil {
		p.metrics.TotalAdd(produce, metrics.Error, eventsCount)

		return fmt.Errorf("p.buildKafkaMessages: %w", err)
	}

	err = p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		p.metrics.TotalAdd(produce, metrics.Error, eventsCount)

		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	p.metrics.TotalAdd(produce, metrics.Ok, eventsCount)

	return nil
}

func buildKafkaMessages(events []domain.Event) ([]kafka.Message, error) {
	messages := make([]kafka.Message, 0, len(events))

	for _, e := range events {
		topic, err := topicByEventType(e.Type)
		if err != nil {
			return nil, fmt.Errorf("topicByEventType: %w", err)
		}

		headers, err := headersFromTraceJSON(e.TraceContext)
		if err != nil {
			return nil, err
		}

		msg := kafka.Message{
			Topic:   topic,
			Key:     []byte(e.ID.String()), // Ключ гарантирует порядок сообщений по профилю
			Value:   e.Value,
			Headers: headers,
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func headersFromTraceJSON(traceCtx []byte) ([]kafka.Header, error) {
	// Если трейс не записан в event, выходим без ошибки
	if len(traceCtx) == 0 {
		return nil, nil
	}

	carrier := propagation.MapCarrier{}

	err := json.Unmarshal(traceCtx, &carrier)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	headers := make([]kafka.Header, 0, len(carrier))
	for k, v := range carrier {
		headers = append(headers, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	return headers, nil
}

func topicByEventType(eventType domain.EventType) (string, error) {
	switch eventType {
	case domain.ProfileCreated:
		return "profiles.created", nil
	default:
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}
}

func (p *Producer) Close() {
	err := p.writer.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka producer: p.writer.Close")
	}

	log.Info().Msg("kafka producer: closed")
}
