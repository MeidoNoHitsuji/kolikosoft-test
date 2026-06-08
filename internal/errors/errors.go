package errors

import (
	"net/http"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model/response"
)

type ErrorStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorStruct(code int, message string) *ErrorStruct {
	return &ErrorStruct{
		Code:    code,
		Message: message,
	}
}

func (e *ErrorStruct) Error() string {
	return e.Message
}

func (e *ErrorStruct) ToResponse() response.Error {
	return response.Error{
		Message: e.Message,
	}
}

var (
	ErrAccountNotFound = NewErrorStruct(http.StatusNotFound, "аккаунт не найден")
	ErrItemNotFound    = NewErrorStruct(http.StatusNotFound, "предмет не найден")
	ErrInternalServer  = NewErrorStruct(http.StatusInternalServerError, "внутренняя ошибка")
	// В идеале нужно просто каждую ошибку разбирать и перезаписывать на правильный вариант, но я обойдусь таким,
	// чтобы не показывать пользователю излишнюю внутрянку
	ErrValidate = NewErrorStruct(http.StatusBadRequest, "ошибка валидации")
)
