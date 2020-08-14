package handlers

import (
	"context"
	"net/http"

	"example.com/back/auth/pkg/app/response"
)

type tokenGetterValidator interface {
	tokenGetter
	validator
}

type blacklistChecker interface {
	IsTokenBlacklisted(ctx context.Context, userID int64, token string) bool
}

// RefreshHandler описывает структуру хэндлера логина
type RefreshHandler struct {
	tgv tokenGetterValidator
	bc  blacklistChecker
}

// NewRefreshHandler возвращает новый хендлер логина
func NewRefreshHandler(bc blacklistChecker, v tokenGetterValidator) *RefreshHandler {
	return &RefreshHandler{tgv: v, bc: bc}
}

// ServeHTTP обрабатывает запрос
func (h *RefreshHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	rt, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		response.JSONError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := h.tgv.ValidateRefreshToken(rt.Value)
	if !ok {
		response.JSONError(rw, "invalid token", http.StatusUnauthorized)
		return
	}

	t := extractTokenFromHeader(r.Header)
	u, _ := h.tgv.ValidateAccessToken(t)

	blacklisted := h.bc.IsTokenBlacklisted(r.Context(), u.ID, rt.Value)
	if blacklisted {
		response.JSONError(rw, "blacklist", http.StatusUnauthorized)
		return
	}

	accessToken := h.tgv.GetAccessToken(u)
	//переустанавливаем в куку новый токен
	setAccessCookie(rw, accessToken)
	// nolint:errcheck
	response.JSON(rw, map[string]string{"accessToken": accessToken}, http.StatusOK)
}
