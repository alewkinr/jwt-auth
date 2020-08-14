package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

//PasswordChangeHandler - описывает структуру handler для смены пароля пользователем
type PasswordChangeHandler struct {
	userGetter userGetter
	validate   validStruct
}

//NewPasswordChangeHandler - возвращает новый handler для смены пароля пользователем
func NewPasswordChangeHandler(ug userGetter, v validStruct) *PasswordChangeHandler {
	return &PasswordChangeHandler{
		userGetter: ug,
		validate:   v,
	}
}

//ServeHTTP - обрабатывает входящие запросы на смену пароля пользователя
func (h *PasswordChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(mux.Vars(r)["userID"], 10, 64)
	if err != nil {
		response.JSONError(w, response.ErrBadRequest, http.StatusBadRequest)
		return
	}
	var req = struct {
		OldPassword string `json:"oldPassword" validate:"required"`
		NewPassword string `json:"newPassword" validate:"required"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONError(w, response.ErrBadRequest, http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.JSONError(w, response.ErrRequestValidationError, http.StatusBadRequest)
		return
	}
	if !h.checkPassword(req.NewPassword) {
		response.JSONError(w, response.ErrRequestValidationError, http.StatusBadRequest)
		return
	}
	u, err := h.userGetter.GetUserByID(r.Context(), int(userID))
	if err != nil {
		if err == user.ErrNotFound {
			response.JSON(w, nil, http.StatusNoContent)
			return
		}
		log.WithField("userID", userID).Error(errors.Wrap(err, "failed GetUserByID"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)); err != nil {
		response.JSONError(w, response.ErrForbidden, http.StatusForbidden)
		return
	}

	if err := h.userGetter.ChangePasswordByUserID(r.Context(), userID, req.NewPassword); err != nil {
		log.WithField("userID", userID).Error(errors.Wrap(err, "failed ChangePasswordByUserID"))
		response.JSONError(w, response.ErrServerInternal, http.StatusInternalServerError)
		return
	}
}

func (h *PasswordChangeHandler) checkPassword(password string) bool {
	if len([]rune(password)) < 6 {
		return false
	}
	isDigital, isLetter := false, false
	for _, f := range password {
		if isDigital && isLetter {
			return true
		}
		if f >= 48 && f <= 57 {
			isDigital = true
			continue
		}
		if (f >= 65 && f <= 90) || (f >= 97 && f <= 122) {
			isLetter = true
			continue
		}
	}
	return isDigital && isLetter
}
