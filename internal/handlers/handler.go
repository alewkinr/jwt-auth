package handlers

import (
	"context"
	"net/http"
	"time"

	"example.com/back/auth/internal/session"

	"example.com/back/auth/internal/user"
)

const (
	refreshTokenCookieName = "refreshToken"
	accessTokenCookieName  = "accessToken"
)

type validStruct interface {
	Struct(interface{}) error
}

// phoneSessionPOSTGetter геттер сессий из бд
type phoneSessionPOSTGetter interface {
	FindActiveByUserPhone(ctx context.Context, usersPhone string) (*session.Session, error)
	FindLastActiveByUserPhone(ctx context.Context, usersPhone string) (*session.Session, error)
}

// phoneSessionPOSTSaver сеттер сессий в БД
type phoneSessionPOSTSaver interface {
	Save(ctx context.Context, usersPhone string, code int64) (*session.Session, error)
	DeleteLastByPhone(ctx context.Context, usersPhone string) error
}

type generateRandom interface {
	GenerateRandomPassword(length, lenDigits, lenSymbol int) (string, error)
}

// notifier отправщик кодов верификации
type notifier interface {
	POSTV1UsersSendMessageSMS(userPhone, message string) error
	POSTV2SendEmailByUserID(ctx context.Context, userID int64,
		templateName string, tags map[string]string) (*http.Response, error)
}

type userGetter interface {
	GetUserByEmailPassword(ctx context.Context, login, password string) (*user.User, error)
	FindUserByEmail(ctx context.Context, email string) (*user.User, error)
	ChangePasswordByUserID(ctx context.Context, userID int64, newPassword string) error
	GetUserByID(ctx context.Context, id int) (*user.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*user.User, error)
}

func setCookie(rw http.ResponseWriter, name, value string, expires time.Time) {
	http.SetCookie(rw, &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   "example.com",
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	})
}

func setRefreshCookie(rw http.ResponseWriter, value string) {
	setCookie(rw, refreshTokenCookieName, value, time.Now().Add(time.Hour*24*7))
}

func setAccessCookie(rw http.ResponseWriter, value string) {
	setCookie(rw, accessTokenCookieName, value, time.Now().Add(time.Minute*10))
}

func deleteCookie(rw http.ResponseWriter, name string) {
	setCookie(rw, name, "", time.Unix(0, 0))
}
