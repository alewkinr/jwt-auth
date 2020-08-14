package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"example.com/back/auth/internal/crypto"
	"example.com/back/auth/internal/session"
	"example.com/back/auth/pkg/app/response"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/../mocks/$GOFILE -package=mocks_phone_sessions_post

// phoneSessionPOSTValidator валидатор
type phoneSessionPOSTValidator interface {
	Struct(interface{}) error
}

// PhoneSessionsPOSTHandler описывает структуру хэндлера создания сессии пользователя по телефону
type PhoneSessionsPOSTHandler struct {
	sm phoneSessionPOSTGetter
	ss phoneSessionPOSTSaver
	v  phoneSessionPOSTValidator
	n  notifier
}

// NewPhoneSessionsPOSTHandler возвращает новый хендлер создания сессии авторизации
func NewPhoneSessionsPOSTHandler(sm phoneSessionPOSTGetter, ss phoneSessionPOSTSaver, v phoneSessionPOSTValidator, notifier notifier) *PhoneSessionsPOSTHandler {
	return &PhoneSessionsPOSTHandler{sm: sm, ss: ss, v: v, n: notifier}
}

// ServeHTTP обслуживает запрос
func (h *PhoneSessionsPOSTHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	type handlerRequest struct {
		Phone       string `json:"phone" validate:"required,max=16,startswith=+"`
		CountryCode string `json:"countryCode" validate:"required,len=2"`
	}

	var req handlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(rw, response.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := h.v.Struct(&req); err != nil {
		response.JSONError(rw, response.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if !isAllowedByRateLimits(r.Context(), h.sm, req.Phone) {
		rw.Header().Add("Retry-After", session.ExpirationTimeout.String())
		response.JSONError(rw, response.ErrRateLimits, http.StatusTooManyRequests)
		return
	}

	if err := h.ss.DeleteLastByPhone(r.Context(), req.Phone); err != nil {
		log.WithField("userPhone", req.Phone).Error(errors.Wrap(err, "failed to delete session on code request"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	code := generateVerificationCode()
	s, err := h.ss.Save(r.Context(), req.Phone, code)
	if err != nil {
		if err == session.ErrDuplicateKey {
			log.WithField("userPhone", req.Phone).Error(errors.Wrap(err, "found duplicate session id"))
			response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
			return
		}
		log.WithField("userPhone", req.Phone).Error(errors.Wrap(err, "failed to save session to DB"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	messageSMS := fmt.Sprintf("%d - ваш код авторизации", code)
	if err := h.n.POSTV1UsersSendMessageSMS(req.Phone, messageSMS); err != nil {
		log.WithField("usersPhone", req.Phone).
			Error(errors.Wrap(err, "failed to send verification code to user"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	// HTTP 200 OK
	type handlerResponse struct {
		SessionID string `json:"sessionId"`
	}
	resp := &handlerResponse{SessionID: s.SessionID.String()}
	response.JSON(rw, resp, http.StatusOK)
}

// generateVerificationCode получаем SMS-код для верификации
func generateVerificationCode() int64 {
	const (
		low  = 1000
		high = 10000
	)
	cr := crypto.New()
	rnd := rand.New(cr)
	return low + rnd.Int63n(high-low)
}
