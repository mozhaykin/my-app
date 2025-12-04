package kafkaproducer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

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
		// BatchSize:    10,
		// BatchTimeout: 100 * time.Millisecond,
	}

	return &Producer{
		config: c,
		writer: w,
	}
}

func (p *Producer) Produce(ctx context.Context, messages ...kafka.Message) error {
	err := p.writer.WriteMessages(ctx, messages...)
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
