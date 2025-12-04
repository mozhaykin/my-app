package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
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
		messageCount, err := w.usecase.OutboxReadAndProduce(context.Background(), w.config.Limit)
		if err != nil {
			log.Error().Err(err).Msg("outbox kafka worker: read and produce failed")
		}

		log.Info().Int("count", messageCount).Msg("outbox kafka worker: read and produce")

		var duration time.Duration

		// если пришло меньше 10 сообщений, значит их больше нет и надо поспать 10 секунд
		if messageCount < w.config.Limit {
			duration = 10 * time.Second

			log.Info().Msg("outbox kafka worker: sleeping 10s")
		}

		select {
		case <-w.stop:
			break FOR // Метка FOR, чтобы выйти не только из select, а полностью из цикла for
		case <-time.After(duration):
		}
	}

	close(w.done)
}

func (w *OutboxKafkaWorker) Stop() {
	log.Info().Msg("outbox kafka worker: closing")

	close(w.stop)

	<-w.done

	log.Info().Msg("outbox kafka worker: closed")
}
