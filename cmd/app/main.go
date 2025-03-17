package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from k8s!"))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}

	fmt.Println("200 OK! Hello handler called")
}

func probe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := chi.NewRouter()

	router.Get("/amozhaykin/my-app/hello", hello)
	router.Get("/live", probe)
	router.Get("/ready", probe)

	fmt.Println("Starting server on :8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
