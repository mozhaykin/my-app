package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppName       string `envconfig:"APP_NAME" required:"true"`
	AppVersion    string `envconfig:"APP_VERSION" required:"true"`
	Level         string `envconfig:"LOGGER_LEVEL" default:"error"`
	PrettyConsole bool   `envconfig:"LOGGER_PRETTY_CONSOLE" default:"false"`
}

func Init(c Config) {
	zerolog.TimeFieldFormat = time.RFC3339    // Формат времени: "2006-01-02T15:04:05Z07:00"
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // Уровень логирования по умолчанию: Info

	level, err := zerolog.ParseLevel(c.Level) // Динамическое изменение уровня логирования, если парсинг успешен
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	log.Logger = log.With(). // Добавление общих полей в логи
					Caller(). // Добавляет файл и номер строки, откуда вызван лог
					Str("app_name", c.AppName).
					Str("app_version", c.AppVersion).
					Logger()

	if c.PrettyConsole { // Если c.PrettyConsole == true, логи выводятся с подсветкой, без JSON (иначе в JSON)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	}

	log.Info().Msg("Logger initialized")
}
