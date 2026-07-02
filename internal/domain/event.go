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
	Status     Status    `json:"status"`
	Verified   bool      `json:"verified"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	OccurredAt time.Time `json:"occurred_at"`
}

type Event struct {
	ID           uuid.UUID `json:"event_id"`
	Type         EventType `json:"event_type"`
	OccurredAt   time.Time `json:"occurred_at"`
	Value        []byte    `json:"value"`
	TraceContext []byte    `json:"traceContext"`
}

type EventType string

const ProfileCreated EventType = "ProfileCreated"

func EventProfileCreated(p Profile) (Event, error) {
	payload := Payload{
		ID:         p.ID,
		Name:       string(p.Name),
		Age:        int(p.Age),
		Status:     p.Status,
		Verified:   p.Verified,
		Email:      p.Contacts.Email,
		Phone:      p.Contacts.Phone,
		OccurredAt: time.Now().UTC(),
	}

	value, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("json.Marshal EventProfileCreated: %w", err)
	}

	return Event{
		ID:         uuid.New(),
		Type:       ProfileCreated,
		OccurredAt: payload.OccurredAt,
		Value:      value,
	}, nil
}
