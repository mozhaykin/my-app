package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
)

type Profile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var log zerolog.Logger

func main() {
	zerolog.TimeFieldFormat = time.RFC3339                   // Инициализация логгера с форматом времени
	log = zerolog.New(os.Stdout).With().Timestamp().Logger() // вывод логов в stdout с форматом времени

	log.Info().Msg("Starting server on :8080") // Логирование старта сервера

	r := chi.NewRouter()

	r.Get("/live", probe)
	r.Get("/ready", probe)

	r.Get("/amozhaykin/my-app/hello", helloGET)
	r.Get("/amozhaykin/my-app/profile", profileGET)
	r.Post("/amozhaykin/my-app/profile", profilePOST)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server") // Логирование ошибки старта сервера
	}
}

func probe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	log.Info(). // логирование в хендлерах
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("Probe handler called")
}

func helloGET(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello!"))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response in helloGET") // Логирование ошибки записи ответа
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("200 OK! Hello handler called")
}

func profileGET(w http.ResponseWriter, r *http.Request) {
	profile := Profile{
		Name: "Alice",
		Age:  30,
	}

	_, err := w.Write([]byte(fmt.Sprintf("Name: %s, Age: %d", profile.Name, profile.Age)))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response in profileGET")
		return
	}

	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("200 OK! Profile handler called")
}

func profilePOST(w http.ResponseWriter, r *http.Request) {
	var newProfile Profile

	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		log.Error().
			Err(err).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("400 Bad Request in profilePOST")
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Interface("profile", newProfile).
		Msg("201 Created! New profile created")
}
