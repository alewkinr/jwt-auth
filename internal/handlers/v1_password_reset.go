package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

//PasswordResetHandler - описывает структуру handler для сброса пароля пользователя
type PasswordResetHandler struct {
	userGetter userGetter
	notifier   notifier
	validate   validStruct
	random     generateRandom
}

//NewPasswordResetHandler - возвращает новый handler для сброса пароля пользователя
func NewPasswordResetHandler(ug userGetter, n notifier, v validStruct, gr generateRandom) *PasswordResetHandler {
	return &PasswordResetHandler{
		userGetter: ug,
		notifier:   n,
		validate:   v,
		random:     gr,
	}
}

//ServeHTTP - обрабатывает входящие запросы на сброс пароля пользователя
func (h *PasswordResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const templateName = "password_recovery"
	var req = struct {
		Email string `json:"email" validate:"required,email"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, response.ErrBadRequest, http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.JSONError(w, response.ErrRequestValidationError, http.StatusBadRequest)
		return
	}
	u, err := h.userGetter.FindUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == user.ErrNotFound {
			response.JSON(w, nil, http.StatusNoContent)
			return
		}
		log.WithField("email", req.Email).Error(errors.Wrap(err, "failed FindUserByEmail"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	pass, err := h.random.GenerateRandomPassword(8, 2, 2)
	if err != nil {
		log.WithField("email", req.Email).Error(errors.Wrap(err, "failed generate new password"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	tags := map[string]string{
		"name":        u.Name,
		"newPassword": pass,
	}
	res, err := h.notifier.POSTV2SendEmailByUserID(r.Context(), u.ID, templateName, tags)
	if err != nil {
		log.WithField("email", req.Email).WithField("userID", u.ID).
			Error(errors.Wrap(err, "failed POSTV2SendEmailByUserID"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		log.WithField("email", req.Email).WithField("userID", u.ID).
			Error(errors.Wrap(err, "not success status response"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	if err := h.userGetter.ChangePasswordByUserID(r.Context(), u.ID, pass); err != nil {
		log.WithField("userID", u.ID).Error(errors.Wrap(err, "failed ChangePasswordByUserID"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
}
