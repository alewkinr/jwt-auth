package session

import (
	"database/sql/driver"
	"time"
)

const (
	ExpirationTimeout = time.Second * 180
	RateLimitTimeout  = time.Second * 30
)

// sessionID id сессии
type sessionID string

// Value реализует интерфейс Valuer
func (s sessionID) Value() (driver.Value, error) {
	return s.String(), nil
}

// String строковое представление роли
func (s sessionID) String() string {
	return string(s)
}

// Scan реализует интерфейс Scanner
func (s *sessionID) Scan(src interface{}) error {
	*s = sessionID(src.(string))
	return nil
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler
func (s *sessionID) UnmarshalJSON(b []byte) error {
	session := string(b)
	*s = sessionID(session)
	return nil
}

// Session объект сессии пользователя
type Session struct {
	SessionID        sessionID `db:"session_id"`
	UsersPhone       string    `db:"users_phone"`
	VerificationCode int64     `db:"verification_code"`
	ExpiresAt        time.Time `db:"expires_at"`
	CreatedAt        time.Time `db:"created_at"`
}
