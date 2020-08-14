package handlers

import (
	"net/http"
	"strconv"

	"example.com/back/auth/internal/permissions"

	"example.com/back/auth/internal/user"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"example.com/back/auth/pkg/app/response"
)

// GetUserHandler описывает структуру хэндлера создания пользователя
type GetUserHandler struct {
	uc         userManager
	permission permissions.IRoleOrOwnerChecker
}

// GetUserResponse структура для ответов данного handler-а
type GetUserResponse struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Name   string `json:"name,omitempty"`
}

// NewGetUserHandler возвращает новый хендлер создания пользователя
func NewGetUserHandler(uc userManager, p permissions.IRoleOrOwnerChecker) *GetUserHandler {
	return &GetUserHandler{
		uc:         uc,
		permission: p,
	}
}

// ServeHTTP обрабатывает запрос
func (h *GetUserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	quID, _ := strconv.Atoi(mux.Vars(r)["userID"])

	if !h.permission.Check(r, int64(quID)) {
		response.JSONError(rw, "Action is not allowed", http.StatusForbidden)
		return
	}
	u, err := h.uc.GetUserByID(r.Context(), quID)
	if err != nil {
		if err == user.ErrNotFound {
			response.JSON(rw, nil, http.StatusNoContent)
			return
		}
		log.Error(errors.Wrap(err, "GetHandler. fetching from DB error"))
		response.JSONError(rw, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.JSON(rw, &GetUserResponse{
		UserID: u.ID,
		Email:  u.Email,
		Phone:  u.Phone,
		Name:   u.Name,
	}, http.StatusOK)
}
