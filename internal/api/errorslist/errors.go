package errorslist

import "errors"

var (
	ErrInternalMsg  = errors.New("Серверная ошибка")
	ErrAccessDenied = errors.New("Доступ запрещен")
)
