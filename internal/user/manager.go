package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const pgErrCodeDuplicateKey = "23505"

var (
	// ErrNotFound запись не найдена
	ErrNotFound = errors.New("not found")
	// ErrInvalidPassword ошибка валидации пароля
	ErrInvalidPassword = errors.New("invalid password")

	// ErrDuplicateKey ошибка дублирования уникального ключа
	ErrDuplicateKey = errors.New("duplicate key")
)

// Manager описывает структуру менеджера по работе с пользователями
type Manager struct {
	db *sqlx.DB
}

// NewManager возвращает новый инстанс менеджера
func NewManager(db *sqlx.DB) *Manager {
	return &Manager{db: db}
}

// GetUserByEmailPassword возвращает пользователя по логину/паролю
func (m *Manager) GetUserByEmailPassword(ctx context.Context, email, password string) (*User, error) {
	u := new(User)
	err := m.db.GetContext(ctx, u, `SELECT * FROM public.user WHERE email = $1`, email)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if berr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); berr != nil {
		return nil, ErrInvalidPassword
	}

	return u, err
}

// GetUserByID возвращает объект пользователя по его ID
func (m *Manager) GetUserByID(ctx context.Context, id int) (*User, error) {
	u := new(User)
	err := m.db.GetContext(ctx, u, `SELECT * FROM public.user WHERE id = $1`, id)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return u, err
}

// GetUserByPhone возвращает объект пользователя по его номеру телефона
func (m *Manager) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	u := new(User)
	err := m.db.QueryRowx(`SELECT * FROM public.user WHERE phone = $1`, phone).
		StructScan(u)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return u, err
}

// GetUserByEmail - поиск пользователя по email
func (m *Manager) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	u := new(User)
	err := m.db.QueryRowxContext(ctx, `SELECT id, name, email, password, phone, 
											role, status, created_on, last_login FROM public.user
											WHERE email = $1`, email).StructScan(u)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return u, err
}

//ChangePassword - смена пароля пользователя по id
func (m *Manager) ChangePasswordByUserID(ctx context.Context, userID int64, newPassword string) error {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "failed bcrypt.GenerateFromPassword")
	}
	_, err = m.db.ExecContext(ctx, `UPDATE public.user SET password = $1 WHERE id = $2`,
		string(passwordBytes), userID)
	if err != nil {
		return errors.Wrap(err, "failed UPDATE password")
	}
	return nil
}

// Create создает пользователя
func (m *Manager) Create(ctx context.Context, name, email, password, phone string, role Role) (*User, error) {
	u := &User{
		Name:      name,
		Email:     email,
		Password:  password,
		Phone:     phone,
		Role:      role,
		Status:    StatusNeedsActivation,
		CreatedOn: time.Now(),
	}

	if u.Password != "" {
		passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		u.Password = string(passwordBytes)
	}
	q := `INSERT INTO public.user 
			(name, email, password, phone, role, status, created_on) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`
	var id int64
	err := m.db.GetContext(ctx, &id, q, u.Name, u.Email, u.Password, u.Phone, u.Role, u.Status, u.CreatedOn)

	// check if returned error is about duplicated errors
	if err, ok := err.(*pgconn.PgError); ok {
		if err.Code == pgErrCodeDuplicateKey {
			return nil, ErrDuplicateKey
		}
	}
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u, err
}

// BlacklistUserToken отправляет токен пользователя в блеклист
func (m *Manager) BlacklistUserToken(ctx context.Context, userID int64, token string) error {
	_, err := m.db.ExecContext(ctx, "INSERT INTO blacklist (user_id, token) VALUES ($1, $2)", userID, token)
	return err
}

// IsTokenBlacklisted проверяет, находится ли токен в блеклисте
func (m *Manager) IsTokenBlacklisted(ctx context.Context, userID int64, token string) bool {
	var i int64
	err := m.db.GetContext(ctx, &i, `SELECT user_id FROM blacklist WHERE user_id = $1 AND token = $2`, userID, token)
	return err != sql.ErrNoRows
}
