package tools

import (
	"fmt"
	rkerror "github.com/rookie-ninja/rk-entry/v2/error"
	"net/http"
)
//==================== 自定义 重写官方rkerror.ErrorInterface接口=====================
type ErrorInterface interface {
	Error() string

	Code() int

	Message() string

	Details() []interface{}
}

type ErrorBuilder interface {
	New(code int, msg string, details ...interface{}) ErrorInterface

	NewCustom() ErrorInterface
}
//======================== 官方Demo =================
type MyError struct {
	ErrCode    int
	ErrMsg     string
	ErrDetails []interface{}
}

func (m MyError) Error() string {
	return fmt.Sprintf("%d-%s", m.ErrCode, m.ErrMsg)
}

func (m MyError) Code() int {
	return m.ErrCode
}

func (m MyError) Message() string {
	return m.ErrMsg
}

func (m MyError) Details() []interface{} {
	return m.ErrDetails
}
type MyErrorBuilder struct{}

func (m *MyErrorBuilder) New(code int, msg string, details ...interface{}) rkerror.ErrorInterface {
	return &MyError{
		ErrCode:    code,
		ErrMsg:     msg,
		ErrDetails: details,
	}
}

func (m *MyErrorBuilder) NewCustom() rkerror.ErrorInterface {
	return &MyError{
		ErrCode:    http.StatusInternalServerError,
		ErrMsg:     "Internal Error",
		ErrDetails: []interface{}{},
	}
}