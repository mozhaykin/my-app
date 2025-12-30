package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
)

func TestNewStatus(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected domain.Status
	}{
		{"pending", "pending", domain.Pending},
		{"active", "active", domain.Active},
		{"inactive", "inactive", domain.Inactive},
		{"banned", "banned", domain.Banned},
		{"unknown", "something_else", domain.Unknown},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			status := domain.NewStatus(c.input)
			require.Equal(t, c.expected, status)
		})
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		name     string
		status   domain.Status
		expected string
	}{
		{"pending", domain.Pending, "pending"},
		{"active", domain.Active, "active"},
		{"inactive", domain.Inactive, "inactive"},
		{"banned", domain.Banned, "banned"},
		{"unknown", domain.Unknown, "unknown"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.expected, c.status.String())
		})
	}
}
