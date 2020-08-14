package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"example.com/back/auth/internal/session"
	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

type PhoneSessionsPUTHandler struct {
	sm phoneSessionPUTGetter
	ss phoneSessionPUTSetter
	v  phoneSessionPUTValidator
	uc userManager
	tm phoneSessionPUTTokenGetter
}

// phoneSessionPOSTValidator валидатор
type phoneSessionPUTValidator interface {
	Struct(interface{}) error
}

// phoneSessionPOSTGetter геттер сессий из бд
type phoneSessionPUTGetter interface {
	FindActiveBySessionID(ctx context.Context, sessionID string) (*session.Session, error)
}

// phoneSessionPUTSetter сеттер сессий в БД
type phoneSessionPUTSetter interface {
	Delete(ctx context.Context, sessionID string) error
	SessionVerified(ctx context.Context, sessionID string) error
}

//phoneSessionPUTTokenGetter геттер доступа к токен
type phoneSessionPUTTokenGetter interface {
	GetAccessToken(u *user.User) string
	GetRefreshToken(u *user.User) string
}

// NewPhoneSessionsPOSTHandler возвращает новый хендлер подтверждения сессии авторизации
func NewPhoneSessionsPUTHandler(
	sm phoneSessionPUTGetter,
	ss phoneSessionPUTSetter,
	uc userManager,
	tm phoneSessionPUTTokenGetter,
	v phoneSessionPUTValidator) *PhoneSessionsPUTHandler {
	return &PhoneSessionsPUTHandler{sm: sm, ss: ss, uc: uc, tm: tm, v: v}
}

// ServeHTTP обслуживает запрос
// nolint: funlen
func (h *PhoneSessionsPUTHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	sessionID := param["sessionId"]
	if len([]rune(sessionID)) < 36 {
		response.JSONError(rw, response.ErrBadRequest, http.StatusBadRequest)
		return
	}

	type handlerRequest struct {
		Phone string `json:"phone" validate:"required,max=16,startswith=+"`
		Code  int64  `json:"code" validate:"required"`
	}
	req := &handlerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(rw, response.ErrBadRequest, http.StatusBadRequest)
		return
	}

	err := h.v.Struct(req)
	if err != nil {
		log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "not valid request"))
		response.JSONError(rw, response.ErrBadRequest, http.StatusBadRequest)
		return
	}
	//поиск активной сессии
	s, err := h.sm.FindActiveBySessionID(r.Context(), sessionID)
	if err != nil {
		if err == session.ErrRecordNotFound {
			log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "not found session id"))
			response.JSONError(rw, response.ErrActiveSessionNotFound, http.StatusNoContent)
			return
		}
		log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "failed to find session to DB"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	// Проверяем рейт-лимит по времени.
	if !h.isAllowedByRateLimits(s) {
		rw.Header().Add("Retry-After", session.ExpirationTimeout.String())
		response.JSONError(rw, response.ErrRateLimits, http.StatusTooManyRequests)
		return
	}

	// проверка телефона, что телефон указанный в запросе совпадает с телефоном в активной сессии
	if !h.isPhoneValid(s.UsersPhone, req.Phone) {
		log.WithField("userSessionID", req.Phone).Error(errors.Wrap(err, "the phone number in the session does not match the phone number from the request"))
		response.JSONError(rw, response.ErrPhoneNotMatchToRequest, http.StatusNoContent)
		return
	}
	// проверка кода подтвержденя, что код указанный в запросе совпадает с кодом в активной сессии
	if !h.isCodeValid(s.VerificationCode, req.Code) {
		response.JSONError(rw, response.ErrCodeNotMatchToRequest, http.StatusNoContent)
		return
	}

	u, err := h.getUserByPhone(r.Context(), req.Phone)
	if err != nil {
		log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "error accurred while get user by phone"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	// указанного телефона нет в example.com, это новый клиент.
	if u == nil {
		if err := h.ss.SessionVerified(r.Context(), sessionID); err != nil {
			log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "failed to set verified session to DB"))
			response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
			return
		}
		response.JSON(rw, map[string]string{"sessionId": sessionID}, http.StatusOK)
		return
	}

	//Указанный телефонный номер есть в public.users, это вернувшийся клиент.
	if err := h.ss.SessionVerified(r.Context(), sessionID); err != nil {
		log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "failed to set verified session to DB"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	//Удаляем сессию из public.sessions
	if err := h.ss.Delete(r.Context(), sessionID); err != nil {
		log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "failed to delete session to DB"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}

	//Получаем токены для пользователя
	accessToken, refreshToken, err := h.getTokens(u)
	if err != nil {
		log.WithField("userSessionID", sessionID).Error(errors.Wrap(err, "failed to get access or refresh token"))
		response.JSONError(rw, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	setRefreshCookie(rw, refreshToken)
	setAccessCookie(rw, accessToken)
	res := struct {
		UserID       int64  `json:"userId"`
		SessionID    string `json:"sessionId"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		UserID:       u.ID,
		SessionID:    sessionID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.JSON(rw, res, http.StatusOK)
}

func (h *PhoneSessionsPUTHandler) getTokens(u *user.User) (accessToken, refreshToken string, err error) {
	accessToken = h.tm.GetAccessToken(u)
	refreshToken = h.tm.GetRefreshToken(u)
	if strings.TrimSpace(accessToken) == "" {
		return "", "", errors.New("an error accurred while get accessToken, token is empty")
	}
	if strings.TrimSpace(refreshToken) == "" {
		return "", "", errors.New("an error accurred while get refreshToken, token is empty")
	}
	return accessToken, refreshToken, nil
}

func (h *PhoneSessionsPUTHandler) getUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	u, err := h.uc.GetUserByPhone(ctx, phone)
	if err != nil {
		//пользователь не найден, возвращаем все nil
		if err == user.ErrNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to find user to DB")
	}

	return u, nil
}

func (h *PhoneSessionsPUTHandler) isPhoneValid(wantedPhone, expectedPhone string) bool {
	return expectedPhone == wantedPhone
}

func (h *PhoneSessionsPUTHandler) isCodeValid(wantedCode, expectedCode int64) bool {
	return expectedCode == wantedCode
}

// isAllowedByRateLimits проверяет возможности клиента продолжить работать с сессиями
func (h *PhoneSessionsPUTHandler) isAllowedByRateLimits(s *session.Session) bool {
	now := time.Now().UTC()
	return now.Before(s.ExpiresAt.UTC())
}
