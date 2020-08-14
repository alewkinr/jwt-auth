package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"example.com/back/auth/internal/user"
)

func TestManager_Validate(t *testing.T) {
	secret := "secret"
	u := &user.User{
		ID:   42,
		Role: user.RoleUser,
	}

	t.Run("invalid_token_signature", func(t *testing.T) {
		m := NewManager(secret, "")

		token := m.GetAccessToken(u)
		token += "foo"

		_, ok := m.ValidateAccessToken(token)
		require.False(t, ok)
	})

	t.Run("empty token", func(t *testing.T) {
		m := NewManager(secret, "")

		token := ""

		_, ok := m.ValidateAccessToken(token)
		require.False(t, ok)
	})

	t.Run("token_without_secret", func(t *testing.T) {
		fakeManager := Manager{accessSecret: nil}
		fakeToken := fakeManager.GetAccessToken(u)

		m := NewManager(secret, "")

		_, ok := m.ValidateAccessToken(fakeToken)
		require.False(t, ok)
	})
}
func TestManager_AdminToken(t *testing.T) {
	secret := ""
	u := &user.User{
		ID:     1,
		Role:   user.RoleAdmin,
		Status: user.StatusActive,
	}

	dur := time.Hour * 24 * 365 * 5 // 5 лет

	m := NewManager(secret, "")
	tkn := m.getAccessToken(u, dur)
	t.Log(tkn)
}
