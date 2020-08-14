package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

type validator interface {
	ValidateAccessToken(token string) (*user.User, bool)
	ValidateRefreshToken(token string) (*user.User, bool)
}

// ValidateHandler описывает структуру хэндлера логина
type ValidateHandler struct {
	v validator
}

// NewValidateHandler возвращает новый хендлер логина
func NewValidateHandler(v validator) *ValidateHandler {
	return &ValidateHandler{v: v}
}

// ServeHTTP обрабатывает запрос
func (h *ValidateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	t := extractTokenFromHeader(r.Header)

	u, ok := h.v.ValidateAccessToken(t)
	if !ok {
		response.JSONError(rw, "invalid token", http.StatusUnauthorized)
		return
	}

	rw.Header().Set("X-Zig-User-ID", fmt.Sprintf("%d", u.ID))
	rw.Header().Set("X-Zig-User-Role", u.Role.String())
	rw.Header().Set("X-Zig-User-Status", u.Status.String())

	rw.WriteHeader(http.StatusOK)
}

func extractTokenFromHeader(header http.Header) string {
	h := header.Get("Authorization")

	return strings.Replace(h, "Bearer ", "", 1)
}
