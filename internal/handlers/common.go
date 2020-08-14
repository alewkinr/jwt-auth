package handlers

import (
	"context"
	"time"

	"example.com/back/auth/internal/session"

	log "github.com/sirupsen/logrus"
)

// isAllowedByRateLimits проверяет возможности клиента продолжить работать с сессиями
func isAllowedByRateLimits(ctx context.Context, sessionGetter phoneSessionPOSTGetter, usersPhone string) bool {
	s, err := sessionGetter.FindLastActiveByUserPhone(ctx, usersPhone)
	if err != nil {
		switch err {
		// сессий для указанного телефона нет
		case session.ErrRecordNotFound:
			return true
		default:
			log.WithField("usersPhone", usersPhone).Error("failed to check rate limits for user")
			return false
		}
	}
	now := time.Now().UTC()
	return now.Sub(s.CreatedAt.UTC()) >= session.RateLimitTimeout
}
