package worden_test

import (
	"runtime"
)

// 获取正在运行的函数名
func RunFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}
