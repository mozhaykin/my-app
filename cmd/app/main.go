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

func main() {
	router := chi.NewRouter()
	router.Get("/amozhaykin/my-app/hello", hello)

	fmt.Println("Starting server on :8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
