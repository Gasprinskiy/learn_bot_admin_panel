package global

import (
	"errors"
	"net/http"
)

var (
	// ErrNoData данные не найдены"
	ErrNoData = errors.New("данные не найдены")
	// ErrInternalError внутряя ошибка
	ErrInternalError = errors.New("произошла внутреняя ошибка")
	// ErrPermissionDenied отказано в доступе
	ErrPermissionDenied = errors.New("отказано в доступе")
	// ErrInvalidParam не верные параметры
	ErrInvalidParam = errors.New("не верные параметры")
	// ErrExpired время вышло
	ErrExpired = errors.New("время вышло")
)

var ErrStatusCodes = map[error]int{
	ErrNoData:        http.StatusNotFound,
	ErrInternalError: http.StatusInternalServerError,
	// ErrInvalidParam:           http.StatusBadRequest,
	// ErrInvalidLoginOrPassword: http.StatusUnauthorized,
	// ErrUserAllreadyExists:     http.StatusConflict,
	// ErrNotAllowedToUse:        http.StatusUnauthorized,
	// ErrExpiredSesstion:        http.StatusUnauthorized,
}

var SuccessStatuses = map[string]int{
	// SuccessLogedOut: http.StatusOK,
}
