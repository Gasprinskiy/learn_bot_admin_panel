package global

import (
	"errors"
	"net/http"
)

var (
	// ErrNoData данные не найдены"
	ErrNoData = errors.New("NoData")
	// ErrInternalError внутряя ошибка
	ErrInternalError = errors.New("InternalError")
	// ErrPermissionDenied отказано в доступе
	ErrPermissionDenied = errors.New("PermissionDenied")
	// ErrInvalidParam не верные параметры
	ErrInvalidParam = errors.New("InvalidParam")
	// ErrExpired время вышло
	ErrExpired = errors.New("Expired")
)

var ErrStatusCodes = map[error]int{
	ErrNoData:        http.StatusNotFound,
	ErrInternalError: http.StatusInternalServerError,
	ErrExpired:       http.StatusGone,
	// ErrInvalidParam:           http.StatusBadRequest,
	// ErrInvalidLoginOrPassword: http.StatusUnauthorized,
	// ErrUserAllreadyExists:     http.StatusConflict,
	// ErrNotAllowedToUse:        http.StatusUnauthorized,
	// ErrExpiredSesstion:        http.StatusUnauthorized,
}

var SuccessStatuses = map[string]int{
	// SuccessLogedOut: http.StatusOK,
}
