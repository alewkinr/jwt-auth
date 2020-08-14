package handlers

import (
	"context"
	"net/http"

	"example.com/back/auth/pkg/app/response"
)

type tokenBlacklister interface {
	BlacklistUserToken(ctx context.Context, userID int64, token string) error
}

// LogoutHandler описывает структуру хэндлера логаута
type LogoutHandler struct {
	tb  tokenBlacklister
	tgv tokenGetterValidator
}

// NewLogoutHandler возвращает новый хендлер логаута
func NewLogoutHandler(tb tokenBlacklister, tgv tokenGetterValidator) *LogoutHandler {
	return &LogoutHandler{tb: tb, tgv: tgv}
}

// ServeHTTP обрабатывает запрос
func (h *LogoutHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rt, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		response.JSONError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	u, ok := h.tgv.ValidateRefreshToken(rt.Value)
	if !ok {
		response.JSONError(rw, "invalid token", http.StatusUnauthorized)
		return
	}

	err = h.tb.BlacklistUserToken(r.Context(), u.ID, rt.Value)
	if err != nil {
		response.JSONError(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	deleteCookie(rw, refreshTokenCookieName)
	deleteCookie(rw, accessTokenCookieName)
	rw.Header().Set("Content-type", "application/json")

	response.JSON(rw, map[string]bool{"success": true}, http.StatusOK)
}
