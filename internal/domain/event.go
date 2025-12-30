package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Age        int       `json:"age"`
	Status     string    `json:"status"`
	Verified   bool      `json:"verified"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	OccurredAt time.Time `json:"occurred_at"`
}

type Event struct {
	ID           uuid.UUID `json:"event_id"`
	Type         EventType `json:"event_type"`
	OccurredAt   time.Time `json:"occurred_at"`
	Value        []byte    `json:"value"` // сериализованный JSON события
	TraceContext []byte    `json:"traceContext"`
}

type EventType string

const ProfileCreated EventType = "ProfileCreated"

func EventProfileCreated(p Profile) (Event, error) {
	// Создаём чистое domain-событие
	payload := Payload{
		ID:         p.ID,
		Name:       string(p.Name),
		Age:        int(p.Age),
		Status:     p.Status.String(),
		Verified:   p.Verified,
		Email:      p.Contacts.Email,
		Phone:      p.Contacts.Phone,
		OccurredAt: time.Now().UTC(),
	}

	// Сериализуем payload в JSON
	value, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("json.Marshal EventProfileCreated: %w", err)
	}

	// Возвращаем объект готовый к записи в outbox
	return Event{
		ID:         uuid.New(),
		Type:       ProfileCreated,
		OccurredAt: payload.OccurredAt,
		Value:      value,
	}, nil
}
