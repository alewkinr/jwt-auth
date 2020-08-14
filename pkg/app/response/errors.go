package response

var (
	// ErrServerInternal ошибка в ответ на 500 ответы
	ErrServerInternal = "internal server error"
	// ErrBadRequest ошибка парсинга входных параметров
	ErrBadRequest = "bad request"
	// ErrRequestValidationError ошибка валидации вхолящего запроса. Детали в логах
	ErrRequestValidationError = "request validation failed"
	// ErrForbidden ошибка прав доступа, юзер не имеет доступа к ресурсу
	ErrForbidden = "insufficient permissions"
	// ErrRateLimits ошибка ограничения rate-лимитов
	ErrRateLimits = "too many requests, try again later"
	// ErrActiveSessionNotFound ошибка не найдена активная сессия
	ErrActiveSessionNotFound = "active session not found"
	// ErrPhoneNotMatchToRequest ошибка телефонный номер из входящего запроса не совпадает с номером телефона в
	// акстивной сессии
	ErrPhoneNotMatchToRequest = "bad phone number"
	// ErrCodeNotMatchToRequest ошибка смс код из входящего запроса не совпадает с кодом в акстивной сессии
	ErrCodeNotMatchToRequest = "bad code number"
	// ErrEmailNotFound - email не найден в БД
	ErrEmailNotFound = "email not found"
)
