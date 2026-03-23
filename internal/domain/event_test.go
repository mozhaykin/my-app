package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/mozhaykin/my-app/internal/domain"
)

func TestEventProfileCreated(t *testing.T) {
	// Пример данных для profile
	p := domain.Profile{
		ID:       uuid.New(),
		Name:     "John Doe",
		Age:      30,
		Status:   domain.Pending,
		Verified: false,
		Contacts: domain.Contacts{
			Email: "john@example.com",
			Phone: "+1234567890",
		},
	}

	// Вызов основной функции
	event, err := domain.EventProfileCreated(p)
	require.NoError(t, err)

	// Проверка
	require.NotEqual(t, uuid.Nil, event.ID, "Event ID should not be nil")
	require.Equal(t, domain.ProfileCreated, event.Type)
	require.WithinDuration(t, time.Now().UTC(), event.OccurredAt, time.Second, "OccurredAt should be near now")

	// Проверяем payload JSON
	var payload domain.Payload
	err = json.Unmarshal(event.Value, &payload)
	require.NoError(t, err)

	require.Equal(t, p.ID, payload.ID)
	require.Equal(t, string(p.Name), payload.Name)
	require.Equal(t, int(p.Age), payload.Age)
	require.Equal(t, p.Status.String(), payload.Status)
	require.Equal(t, p.Verified, payload.Verified)
	require.Equal(t, p.Contacts.Email, payload.Email)
	require.Equal(t, p.Contacts.Phone, payload.Phone)
	require.WithinDuration(t, time.Now().UTC(), payload.OccurredAt, time.Second)
}
