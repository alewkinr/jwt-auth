package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/../mocks/mocks_v1_login/$GOFILE -package=mocks_v1_login

type tokenGetter interface {
	GetAccessToken(u *user.User) string
	GetRefreshToken(u *user.User) string
}

// LoginHandler описывает структуру хэндлера логина
type LoginHandler struct {
	ug userGetter
	tg tokenGetter
}

// NewLoginHandler возвращает новый хендлер логина
func NewLoginHandler(ug userGetter, tg tokenGetter) *LoginHandler {
	return &LoginHandler{ug: ug, tg: tg}
}

// ServeHTTP обрабатывает запрос
func (h *LoginHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// nolint:errcheck
	r.ParseForm()
	login := r.FormValue("login")
	password := r.FormValue("password")
	u, err := h.ug.GetUserByEmailPassword(r.Context(), login, password)

	rw.Header().Set("Content-type", "application/json")
	if err != nil {
		log.Error(errors.Wrap(err, "login GetUserByEmailPassword error "))
		response.JSONError(rw, "unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken := h.tg.GetAccessToken(u)
	refreshToken := h.tg.GetRefreshToken(u)

	setAccessCookie(rw, accessToken)
	setRefreshCookie(rw, refreshToken)

	response.JSON(rw, map[string]string{"accessToken": accessToken}, http.StatusOK)
}
