package session

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const errCodeDuplicateKey = "23505"

// Manager описывает структуру менеджера по работе с пользователями
type Manager struct {
	store *sqlx.DB
}

// NewManager возвращает новый инстанс менеджера
func NewManager(db *sqlx.DB) *Manager {
	return &Manager{store: db}
}

// FindActiveBySessionID возвращает объект активной но НЕ ВЕРИФИЦИРОВАННОЙ сессии по sessionID
func (m *Manager) FindActiveBySessionID(ctx context.Context, sessionID string) (*Session, error) {
	q := `SELECT 
				session_id,
				users_phone,
				verification_code,
				expires_at,
				created_at
		FROM
				public.sessions
		WHERE session_id = $1
		AND expires_at > CURRENT_TIMESTAMP
		AND is_verified = false;`

	var s Session
	if err := m.store.GetContext(ctx, &s, q, sessionID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound
		case sql.ErrConnDone:
			return nil, ErrConnectionClosed
		default:
			return nil, errors.Wrap(err, "failed to find users session by users session id")
		}
	}
	return &s, nil
}

// FindActiveByUserPhone возвращает объект активной сессии по userPhone
func (m *Manager) FindActiveByUserPhone(ctx context.Context, usersPhone string) (*Session, error) {
	q := `SELECT 
				session_id,
				users_phone,
				verification_code,
				expires_at,
				created_at
		FROM
				public.sessions
		WHERE users_phone = $1
			AND expires_at > CURRENT_TIMESTAMP;`

	var s Session
	if err := m.store.GetContext(ctx, &s, q, usersPhone); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound
		case sql.ErrConnDone:
			return nil, ErrConnectionClosed
		default:
			return nil, errors.Wrap(err, "failed to find users session by users phone")
		}
	}
	return &s, nil
}

// FindLastActiveByUserPhone возвращает объект последней сессии по телефону пользователя
func (m *Manager) FindLastActiveByUserPhone(ctx context.Context, usersPhone string) (*Session, error) {
	q := `SELECT 
				session_id,
				users_phone,
				verification_code,
				expires_at,
				created_at
		FROM
				public.sessions
		WHERE users_phone = $1
			  AND expires_at >= CURRENT_TIMESTAMP 
		ORDER BY created_at DESC;`
	var s Session
	if err := m.store.QueryRowxContext(ctx, q, usersPhone).StructScan(&s); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound
		case sql.ErrConnDone:
			return nil, ErrConnectionClosed
		default:
			return nil, errors.Wrap(err, "failed to find last users session by users phone")
		}
	}
	return &s, nil
}

// Save сохраняет новую сессию в БД
func (m *Manager) Save(ctx context.Context, usersPhone string, code int64) (*Session, error) {
	q := `INSERT INTO public.sessions (
				session_id,
				users_phone, 
				verification_code, 
				expires_at, 
				created_at
			)
			VALUES (
				:session_id,
				:users_phone,
				:verification_code,
				:expires_at,
				:created_at
			)
			RETURNING session_id;`

	now := time.Now().UTC()
	s := Session{
		SessionID:        sessionID(uuid.New().String()),
		UsersPhone:       usersPhone,
		VerificationCode: code,
		ExpiresAt:        now.Add(ExpirationTimeout),
		CreatedAt:        now,
	}
	if _, err := m.store.NamedQueryContext(ctx, q, s); err != nil {
		switch err.(*pgconn.PgError).Code {
		case errCodeDuplicateKey:
			return nil, ErrDuplicateKey
		default:
			return nil, errors.Wrap(err, "failed to save session to DB")
		}
	}
	return &s, nil
}

// Delete удаляет сессию из БД
func (m *Manager) Delete(ctx context.Context, sessionID string) error {
	q := `DELETE FROM public.sessions
	WHERE session_id = $1;`

	if _, err := m.store.ExecContext(ctx, q, sessionID); err != nil {
		switch err.(*pgconn.PgError).Code {
		case errCodeDuplicateKey:
			return ErrDuplicateKey
		default:
			return errors.Wrap(err, "failed to delete session to DB")
		}
	}
	return nil
}

// DeleteLastByPhone — удаляет сессию по номеру телефона клиента. Нужно чтобы чистило предыдущий код
func (m *Manager) DeleteLastByPhone(ctx context.Context, usersPhone string) error {
	q := `DELETE FROM public.sessions
		WHERE users_phone = $1
		AND is_verified = false
		AND expires_at >= CURRENT_TIMESTAMP`

	if _, err := m.store.ExecContext(ctx, q, usersPhone); err != nil {
		return errors.Wrap(err, "failed to delete session by users phone in DB")
	}
	return nil
}

// SessionVerified отмечаем сессию verified = truе
func (m *Manager) SessionVerified(ctx context.Context, sessionID string) error {
	q := `UPDATE public.sessions
	SET is_verified=true 
	WHERE session_id=$1;`
	if _, err := m.store.ExecContext(ctx, q, sessionID); err != nil {
		switch err.(*pgconn.PgError).Code {
		case errCodeDuplicateKey:
			return ErrDuplicateKey
		default:
			return errors.Wrap(err, "failed to delete session to DB")
		}
	}
	return nil
}

// IsSessionValid проверяет статус верификации сессии
// Возвращает true, только если есть хотябы 1 подтвержденная сессия на юзера
func (m *Manager) IsSessionValid(ctx context.Context, sessionID string) bool {
	q := `SELECT 
				count(session_id) AS count_sessions
			FROM public.sessions
			WHERE session_id = $1
				AND is_verified = true
				AND expires_at > CURRENT_TIMESTAMP;`

	var numSessions int
	if err := m.store.QueryRowxContext(ctx, q, sessionID).Scan(&numSessions); err != nil {
		return false
	}
	return numSessions == 1
}
