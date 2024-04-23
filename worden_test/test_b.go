package worden_test

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
知识点：
1：文件绝对路径
*/
func Test_b(name string) {
	fmt.Println("=========================== 当前函数名：", RunFuncName())
	//判断是否绝对路径,拼接成字符串
	bootConfigPath := "/var/log/aa.log"
	bool_b := filepath.IsAbs(bootConfigPath)
	if !bool_b {
		wd, _ := os.Getwd()
		bootConfigPath = filepath.Join(wd, bootConfigPath)
	}
}
