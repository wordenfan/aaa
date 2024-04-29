package tools

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
//全局配置变量
var (
	IlLegalUser = *NewError(3001,"非法用户")
)

type GinError struct {
	ErrCode    codes.Code
	ErrMsg     string
}

func (m GinError) Error() string {
	return fmt.Sprintf("%d-%s", m.ErrCode, m.ErrMsg)
}

func NewError(code codes.Code,msg string) *GinError {
	return &GinError{
		ErrCode:    code,
		ErrMsg:     msg,
	}
}
func transferToGrpcError(err GinError) error {
	return status.Error(err.ErrCode,err.ErrMsg)
}