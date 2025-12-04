package kafkaconsumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
)

type Config struct {
	Addr     []string `envconfig:"KAFKA_CONSUMER_ADDR"     required:"true"`
	Topic    string   `default:"awesome-topic"             envconfig:"KAFKA_CONSUMER_TOPIC"`
	Group    string   `default:"awesome-group"             envconfig:"KAFKA_CONSUMER_GROUP"`
	Disabled bool     `envconfig:"KAFKA_CONSUMER_DISABLED"`
}
type Consumer struct {
	config  Config
	reader  *kafka.Reader
	usecase *usecase.UseCase
	stop    context.CancelFunc
	done    chan struct{}
}

func New(cfg Config, uc *usecase.UseCase) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Addr,
		Topic:          cfg.Topic,
		GroupID:        cfg.Group,
		ErrorLogger:    logger.ErrorLogger(),
		CommitInterval: 100 * time.Millisecond,
	})

	ctx, stop := context.WithCancel(context.Background())

	c := &Consumer{
		config:  cfg,
		reader:  r,
		usecase: uc,
		stop:    stop,
		done:    make(chan struct{}),
	}

	go c.run(ctx)

	return c
}

func (c *Consumer) run(ctx context.Context) {
	log.Info().Msg("kafka consumer: started")

FOR:
	for {
		// Читаем сообщение из kafka
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				log.Info().Msg("kafka consumer: context canceled")

				break FOR
			case errors.Is(err, io.EOF):
				log.Warn().Err(err).Msg("kafka consumer: FetchMassage")

				break FOR
			}

			log.Error().Err(err).Msg("kafka consumer: FetchMessage")
		}

		log.Info().Str("key", string(m.Key)).Msg("kafka consumer: message received")

		// Тут по идее нужно вызывать метод из usecase для обработки сообщения
		// Я же просто печатаю сообщение
		var p domain.Profile

		err = json.Unmarshal(m.Value, &p)
		if err != nil {
			log.Error().Err(err).Msg("kafka consumer: json.Unmarshal")
		}

		fmt.Printf("Topic: %s, Partition: %d, Offset: %d, Key: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))

		fmt.Printf("Profile: %+v\n", p)

		// Коммитим оффсет в consumer group
		err = c.reader.CommitMessages(ctx, m)
		if err != nil {
			log.Error().Err(err).Msg("kafka consumer: CommitMessages")
		}
	}

	close(c.done)
}

func (c *Consumer) Close() {
	log.Info().Msg("kafka consumer: closing")

	c.stop()

	err := c.reader.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka consumer: reader.Close")
	}

	<-c.done

	log.Info().Msg("kafka consumer: closed")
}
