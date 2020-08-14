package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Role тип роли пользователя
// Использовал число на случай разделения доступа, чтобы можно было сравнивать больше/меньше
type Role uint8

const (
	// RoleUnknown дефолтное значение роли пользователя
	RoleUnknown Role = iota
	// RoleUser роль пользователя
	RoleUser
	// RolePsychologist роль психолога
	RolePsychologist
	// RoleManager роль менеджера
	RoleManager
	// RoleAdmin роль администратора
	RoleAdmin
)

const (
	userRole         = "client"
	psychologistRole = "psychologist"
	managerRole      = "manager"
	adminRole        = "admin"
)

var (
	// ErrUnknownRole неизвестная роль
	ErrUnknownRole = errors.New("unknown role")
	// ErrUnknownStatus  неизвестный статус
	ErrUnknownStatus = errors.New("unknown user status")
)

// Value реализует интерфейс Valuer
func (r Role) Value() (driver.Value, error) {
	return r.String(), nil
}

// Scan реализует интерфейс Scanner
func (r *Role) Scan(src interface{}) error {
	switch src.(string) {
	case userRole:
		*r = RoleUser
	case psychologistRole:
		*r = RolePsychologist
	case managerRole:
		*r = RoleManager
	case adminRole:
		*r = RoleAdmin
	default:
		return ErrUnknownRole
	}
	return nil
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler
func (r *Role) UnmarshalJSON(b []byte) error {
	var role string
	_ = json.Unmarshal(b, &role)
	switch role {
	case userRole:
		*r = RoleUser
	case psychologistRole:
		*r = RolePsychologist
	case managerRole:
		*r = RoleManager
	case adminRole:
		*r = RoleAdmin
	default:
		return ErrUnknownRole
	}
	return nil
}

// String строковое представление роли
func (r Role) String() string {
	switch r {
	case RoleUser:
		return userRole
	case RolePsychologist:
		return psychologistRole
	case RoleManager:
		return managerRole
	case RoleAdmin:
		return adminRole
	default:
		return "unknown"
	}
}

// Status описывает статус пользователя
type Status uint8

const (
	// StatusNeedsActivation требует активации
	StatusNeedsActivation Status = iota
	// StatusActive активный
	StatusActive
)
const (
	inactiveStatus = "needs_activation"
	activeStatus   = "active"
)

// Value реализует интерфейс Valuer
func (s Status) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan реализует интерфейс Scanner
func (s *Status) Scan(src interface{}) error {
	switch src.(string) {
	case activeStatus:
		*s = StatusActive
	case inactiveStatus:
		*s = StatusNeedsActivation
	default:
		return ErrUnknownStatus
	}
	return nil
}

// String строковое представление статуса активации пользователя
func (s Status) String() string {
	switch s {
	case StatusNeedsActivation:
		return inactiveStatus
	case StatusActive:
		return activeStatus
	default:
		return "unknown status"
	}
}

// User описывает структуру юзера
type User struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Phone     string     `db:"phone"`
	Role      Role       `db:"role"`
	Status    Status     `db:"status"`
	CreatedOn time.Time  `db:"created_on"`
	LastLogin *time.Time `db:"last_login"`
}
