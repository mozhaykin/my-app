package v1

import (
	"encoding/json"
	"net/http"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h *Handlers) GetProfile(w http.ResponseWriter, _ *http.Request) {
	output := dto.CreateProfileOutput{
		Name: "Alice",
		Age:  30,
	}

	data, err := json.Marshal(&output)
	if err != nil {
		http.Error(w, "json error", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "write error", http.StatusBadRequest)
	}
}
