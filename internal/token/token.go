package token

import (
	"time"

	"github.com/brianvoe/sjwt"

	"example.com/back/auth/internal/user"
)

const (
	accessTokenLifetime  = 1 * 24 * time.Hour       // 24 часа
	refreshTokenLifetime = 7 * 24 * 183 * time.Hour // ~ 6 месяцев
)

// Manager описывает структуру менеджера токенов
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
}

// NewManager возвращает инстанс менеджера токенов
func NewManager(accessSecret, refreshSecret string) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

// GetAccessToken возвращает access token по пользователю
func (m *Manager) GetAccessToken(u *user.User) string {
	return m.getAccessToken(u, accessTokenLifetime)
}

// GetAccessToken возвращает access token по пользователю
func (m *Manager) getAccessToken(u *user.User, lifetime time.Duration) string {
	claims := sjwt.New()

	claims.Set("i", u.ID)
	claims.Set("r", u.Role)
	claims.Set("s", u.Status)

	claims.SetExpiresAt(time.Now().Add(lifetime))

	return claims.Generate(m.accessSecret)
}

// GetRefreshToken возвращает refresh token по пользователю
func (m *Manager) GetRefreshToken(u *user.User) string {
	claims := sjwt.New()

	claims.Set("i", u.ID)
	claims.SetExpiresAt(time.Now().Add(refreshTokenLifetime))

	return claims.Generate(m.refreshSecret)
}

// ValidateAccessToken валидирует access токен и возвращает юзера
func (m *Manager) ValidateAccessToken(token string) (*user.User, bool) {
	return m.validate(token, m.accessSecret)
}

// ValidateRefreshToken валидирует refresh токен и возвращает юзера
func (m *Manager) ValidateRefreshToken(token string) (*user.User, bool) {
	return m.validate(token, m.refreshSecret)
}

// validate валидирует токен и возвращает юзера
func (m *Manager) validate(token string, secret []byte) (*user.User, bool) {
	verified := sjwt.Verify(token, secret)
	claims, err := sjwt.Parse(token)
	u := claimsToUser(claims)
	if !verified || err != nil {
		return u, false
	}

	if err := claims.Validate(); err != nil {
		return u, false
	}

	return u, true
}

func claimsToUser(claims sjwt.Claims) *user.User {
	id, _ := claims.GetInt("i")
	role, _ := claims.GetInt("r")
	status, _ := claims.GetInt("s")

	u := &user.User{
		ID:     int64(id),
		Role:   user.Role(role),
		Status: user.Status(status),
	}

	return u
}
