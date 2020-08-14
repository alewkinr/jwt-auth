package session

import "errors"

var (
	// ErrRecordNotFound стандартизированная ошибка для всех БД
	ErrRecordNotFound = errors.New("record not found")

	// ErrRecorcNotChanged ошибка в случае если после обновления в БД ничего не изменилось
	ErrRecorcNotChanged = errors.New("no changes processed in storage")

	// ErrAccessDenied ошибка в случае, если происходит миссматч userId в заказе и заголовке запроса
	ErrAccessDenied = errors.New("action denied, not enough rights")

	// ErrDuplicateKey ошибка дублирования уникального ключа
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrConnectionClosed ошибка подключения к ДБ
	ErrConnectionClosed = errors.New("DB connection has been terminated")
)
