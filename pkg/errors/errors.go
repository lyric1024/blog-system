package errors

import "fmt"

type ErrCode int

// 数据库连接错误、用户认证失败、文章或评论不存在等
const (
	ErrCodeBadRequest   ErrCode = 400 // 请求参数
	ErrCodeUnauthorized ErrCode = 401 // 未认证
	ErrCodeNotFound     ErrCode = 404 // 不存在
	ErrCodeForbidden    ErrCode = 403 // 禁止
	ErrCodeExsit        ErrCode = 409 // 已存在
	ErrCodeInternal     ErrCode = 500 // 内部错误
	ErrCodeUnknow       ErrCode = 500 // 未知错误
)

type AppError struct {
	Code    ErrCode
	Message string
	Err     error // 原始错误
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("AppError code: %d, message: %s, err: %+v", e.Code, e.Message, e.Err)
	}

	return fmt.Sprintf("AppError code: %d, message: %s", e.Code, e.Message)
}

func NewAppError(code ErrCode, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// ==================================
func BadRequest(msg string, err error) *AppError {
	return NewAppError(ErrCodeBadRequest, msg, err)
}
func Unauthorized(msg string, err error) *AppError {
	return NewAppError(ErrCodeUnauthorized, msg, err)
}
func NotFound(msg string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, msg, err)
}
func Forbidden(msg string, err error) *AppError {
	return NewAppError(ErrCodeForbidden, msg, err)
}
func Exsit(msg string, err error) *AppError {
	return NewAppError(ErrCodeExsit, msg, err)
}
func Internal(msg string, err error) *AppError {
	return NewAppError(ErrCodeInternal, msg, err)
}
func Unknow(msg string, err error) *AppError {
	return NewAppError(ErrCodeUnknow, msg, err)
}
