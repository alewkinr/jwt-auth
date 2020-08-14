package handlers

import (
	"fmt"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v3/is"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app/response"
)

// QueryHandler описывает структуру хэндлера поиска пользователя
type QueryHandler struct {
	uc userManager
}

// NewQueryHandler возвращает новый хендлер поиска пользователя
func NewQueryHandler(uc userManager) *QueryHandler {
	return &QueryHandler{uc: uc}
}

// ServeHTTP обрабатывает запрос
func (h *QueryHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	query := params.Get("phone")
	e164 := fmt.Sprintf("+%s", strings.ReplaceAll(query, " ", ""))

	if err := validatePhone(e164); err != nil {
		log.Error(errors.Wrap(err, "failed to validate phone number"))
		response.JSONError(rw, "Validation error", http.StatusBadRequest)
		return
	}

	u, err := h.uc.GetUserByPhone(r.Context(), e164)
	if err != nil {
		switch err {
		case user.ErrNotFound:
			response.JSON(rw, nil, http.StatusNoContent)
		default:
			log.Error(errors.Wrap(err, "failed to get user from DB"))
			response.JSONError(rw, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	response.JSON(rw, map[string]int64{"userId": u.ID}, http.StatusOK)
}

func validatePhone(phone string) error {
	return validation.Validate(phone, validation.Required, validation.Length(2, 15), is.E164)
}
