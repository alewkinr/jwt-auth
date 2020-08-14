package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/../mocks/mocks_v1_users_post/$GOFILE -package=mocks_v1_users_post

type userManager interface {
	Create(ctx context.Context, name, email, password, phone string, role user.Role) (*user.User, error)
	GetUserByID(ctx context.Context, id int) (*user.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*user.User, error)
}

// sessionsVerificationChecker интерфейс для проверки валидности sessionId
type sessionsVerificationChecker interface {
	IsSessionValid(ctx context.Context, sessionID string) bool
}

// sessionsDeleter интерфейс для удаления сессий из хранилища
type sessionsDeleter interface {
	Delete(ctx context.Context, sessionID string) error
}

type requestValidator interface {
	Struct(interface{}) error
}

// CreateHandler описывает структуру хэндлера создания пользователя
type CreateHandler struct {
	uc userManager
	tm tokenGetter
	v  requestValidator
	sv sessionsVerificationChecker
	sd sessionsDeleter
}

// CreateRequest описывает структуру запроса на создание пользователя
type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
	// TODO после починки фронта вернуть требования numeric,lt=15
	// Phone string    `json:"phone" validate:"required,lt=15"`
	Phone     string    `json:"phone" validate:"required,lt=18"`
	SessionID string    `json:"sessionId" validate:"required,uuid4"`
	Role      user.Role `json:"role" validate:"required,min=1"`
}

type CreateResponse struct {
	UserID       int    `json:"userId,omitempty"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Error        string `json:"error,omitempty"`
}

// NewCreateHandler возвращает новый хендлер создания пользователя
func NewCreateHandler(uc userManager, tm tokenGetter, v requestValidator,
	sv sessionsVerificationChecker, sd sessionsDeleter) *CreateHandler {
	return &CreateHandler{uc: uc, tm: tm, v: v, sv: sv, sd: sd}
}

// ServeHTTP обрабатывает запрос
func (h *CreateHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONError(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(rw, response.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := h.v.Struct(req); err != nil {
		response.JSONError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO удалить после починки фронта.
	// Временное решение для очистки номера телефона на бэке.
	re := regexp.MustCompile(`[^0-9\+]`)
	req.Phone = re.ReplaceAllString(req.Phone, "")
	if len(req.Phone) > 15 {
		response.JSONError(rw, "Phone number too long", http.StatusBadRequest)
		return
	}
	if len(req.Phone) == 0 {
		response.JSONError(rw, "Phone number too short", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(req.Phone, "+") {
		req.Phone = "+" + req.Phone
	}

	if !h.sv.IsSessionValid(r.Context(), req.SessionID) {
		response.JSONError(rw, "sessionId is not valid", http.StatusForbidden)
		return
	}

	// создаем пользователя в БД и возвращаем для него токены
	u, err := h.uc.Create(r.Context(), req.Name, req.Email, "", req.Phone, req.Role)
	if err != nil {
		if err == user.ErrDuplicateKey {
			response.JSONError(rw, "User already exists", http.StatusConflict)
		} else {
			response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		}
		return
	}

	// удаляем сессию, на которую выпустили токены
	if err := h.sd.Delete(r.Context(), req.SessionID); err != nil {
		log.WithField("sessionId", req.SessionID).Error(errors.Wrap(err, "failed to delete session from DB"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	at, rt := h.tm.GetAccessToken(u), h.tm.GetRefreshToken(u)
	//  устанавливаем рефреш и access токен в куку
	setRefreshCookie(rw, rt)
	setAccessCookie(rw, at)
	//  HTTP 200
	response.JSON(rw, &CreateResponse{UserID: int(u.ID), AccessToken: at, RefreshToken: rt}, http.StatusOK)
}
