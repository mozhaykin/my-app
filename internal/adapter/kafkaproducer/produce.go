package kafkaproducer

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
)

type Config struct {
	Addr []string `envconfig:"KAFKA_WRITER_ADDR" required:"true"`
}

type Producer struct {
	config Config
	writer *kafka.Writer
}

func New(c Config) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.Addr...),
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  logger.ErrorLogger(),
		Async:        false, // можно включить отправку батчами (с интервалом в 1s) если producer захлёбывается
	}

	return &Producer{
		config: c,
		writer: w,
	}
}

func (p *Producer) Produce(ctx context.Context, events ...domain.Event) error {
	//nolint:prealloc
	var msgs []kafka.Message

	for _, event := range events {
		msg := kafka.Message{
			Topic: event.Topic,
			Key:   event.Key,
			Value: event.Value,
		}

		msgs = append(msgs, msg)
	}

	err := p.writer.WriteMessages(ctx, msgs...)
	if err != nil {
		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	return nil
}

func (p *Producer) Close() {
	err := p.writer.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka producer: p.writer.Close")
	}

	log.Info().Msg("kafka producer: closed")
}
