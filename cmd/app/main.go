package main

import (
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/app"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/logger"
)

func main() {
	logger.Init(logger.Config{
		AppName:       "my-app",
		AppVersion:    "v0.1.0",
		Level:         "debug",
		PrettyConsole: true,
	})

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
