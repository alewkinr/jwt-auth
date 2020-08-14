package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"example.com/back/auth/internal/session"
	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

//PostV2PhoneSessionsHandler - структура handler для создании сессии пользователя по телефону версии 2
type PostV2PhoneSessionsHandler struct {
	userGetter    userGetter
	sessionGetter phoneSessionPOSTGetter
	sessionSaver  phoneSessionPOSTSaver
	notifier      notifier
	valid         validStruct
}

//NewPostV2PhoneSessionsHandler - инициализация handler для создании сессии пользователя по телефону версии 2
func NewPostV2PhoneSessionsHandler(u userGetter, sg phoneSessionPOSTGetter, ss phoneSessionPOSTSaver,
	n notifier, v validStruct) *PostV2PhoneSessionsHandler {
	return &PostV2PhoneSessionsHandler{
		userGetter:    u,
		sessionGetter: sg,
		sessionSaver:  ss,
		notifier:      n,
		valid:         v,
	}
}

//ServeHTTP - обрабатывает входящие запросы для handler создания сессий пользователя по телефону версии 2
func (h *PostV2PhoneSessionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req = struct {
		Phone       string `json:"phone" validate:"required,max=16,startswith=+"`
		CountryCode string `json:"countryCode" validate:"required,len=2"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, response.ErrBadRequest, http.StatusBadRequest)
		return
	}
	if err := h.valid.Struct(&req); err != nil {
		response.JSONError(w, response.ErrRequestValidationError, http.StatusBadRequest)
		return
	}
	_, err := h.userGetter.GetUserByPhone(r.Context(), req.Phone)
	if err != nil {
		if err == user.ErrNotFound {
			response.JSON(w, nil, http.StatusNoContent)
			return
		}
		log.WithField("phone", req.Phone).Error(errors.Wrap(err, "failed GetUserByPhone"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	if !isAllowedByRateLimits(r.Context(), h.sessionGetter, req.Phone) {
		w.Header().Add("Retry-After", session.ExpirationTimeout.String())
		response.JSONError(w, response.ErrRateLimits, http.StatusTooManyRequests)
		return
	}

	if err := h.sessionSaver.DeleteLastByPhone(r.Context(), req.Phone); err != nil {
		log.WithField("phone", req.Phone).Error(errors.Wrap(err, "failed DeleteLastByPhone"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	code := generateVerificationCode()
	s, err := h.sessionSaver.Save(r.Context(), req.Phone, code)
	if err != nil {
		if err == session.ErrDuplicateKey {
			log.WithField("phone", req.Phone).Error(errors.Wrap(err, "found duplicate session id"))
			response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
			return
		}
		log.WithField("phone", req.Phone).Error(errors.Wrap(err, "failed sessionSaver.Save"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	messageSMS := fmt.Sprintf("%d - ваш код для авторизации", code)
	if err := h.notifier.POSTV1UsersSendMessageSMS(req.Phone, messageSMS); err != nil {
		log.WithField("usersPhone", req.Phone).
			Error(errors.Wrap(err, "failed POSTV1UsersSendMessageSMS"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	response.JSON(w, map[string]string{"sessionId": s.SessionID.String()}, http.StatusOK)
}
