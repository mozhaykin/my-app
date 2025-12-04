package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

func TestNewProfile(t *testing.T) {
	cases := []struct {
		name    string
		age     int
		email   string
		phone   string
		wantErr bool
	}{
		{"Valid Profile", 25, "test@example.com", "+7123456789", false},
		{"In", 25, "test@example.com", "+7123456789", true},
		{"Invalid Age Min", 17, "test@example.com", "+7123456789", true},
		{"Invalid Age Max", 121, "test@example.com", "+7123456789", true},
		{"Invalid Email", 25, "invalid-email", "+7123456789", true},
		{"Invalid Phone", 25, "test@example.com", "invalid-phone", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			profile, err := domain.NewProfile(c.name, c.age, c.email, c.phone)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEqual(t, uuid.Nil, profile.ID)
				require.Equal(t, c.name, string(profile.Name))
				require.Equal(t, c.age, int(profile.Age))
				require.Equal(t, c.email, profile.Contacts.Email)
				require.Equal(t, c.phone, profile.Contacts.Phone)
				require.Equal(t, domain.Pending, profile.Status)
				require.False(t, profile.Verified)
			}
		})
	}
}

func TestProfile_IsDeleted(t *testing.T) {
	var p domain.Profile
	require.False(t, p.IsDeleted())

	p.DeletedAt = time.Now()
	require.True(t, p.IsDeleted())
}

func TestProfile_ToKafkaMsg(t *testing.T) {
	id := uuid.New()
	p := domain.Profile{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      domain.Name("John"),
		Age:       domain.Age(30),
		Status:    domain.Status(1),
		Verified:  true,
		Contacts: domain.Contacts{
			Phone: "+79123456789",
			Email: "test@example.com",
		},
	}

	topic := "profile.created"

	msg, err := p.ToKafkaMsg(topic)
	require.NoError(t, err)

	require.Equal(t, topic, msg.Topic)
	require.Equal(t, []byte(id.String()), msg.Key)

	// Проверяем, что JSON идентичен ожидаемому
	var decoded domain.Profile
	err = json.Unmarshal(msg.Value, &decoded)
	require.NoError(t, err)
	require.Equal(t, p.ID, decoded.ID)
	require.Equal(t, p.Name, decoded.Name)
	require.Equal(t, p.Age, decoded.Age)
	require.Equal(t, p.Status, decoded.Status)
	require.Equal(t, p.Verified, decoded.Verified)
	require.Equal(t, p.Contacts, decoded.Contacts)
	// т.к. во время теста время меняется (на наносекунды),
	// то чтобы тест прошел без ошибок, сравниваем время с округлением до секунды
	require.WithinDuration(t, p.CreatedAt, decoded.CreatedAt, time.Second)
	require.WithinDuration(t, p.UpdatedAt, decoded.UpdatedAt, time.Second)
}
