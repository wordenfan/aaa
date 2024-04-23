package worden_test

import (
	"fmt"
	"golang.org/x/net/context"
)

/*
知识点：
1：map的键值引用会返回2个值：
eg：

	if hook, ok := hooks["category1"]["hookName1"]; ok {
		hook(ctx)
	}

2：Hook的定义引用方式学习，将函数放入数组然后hook调用
*/
func Test_HookFunc(name string) {
	fmt.Println("=========================== 当前函数名：", RunFuncName())
	// 创建一个新的hookFuncM实例
	hooks := newHookFuncM()
	// 先初始化外层map
	if _, ok := hooks["category1"]; !ok {
		hooks["category1"] = make(map[string]func(context.Context))
	}
	// 给hooks添加函数
	hooks["category1"]["hookName1"] = exampleHookFunc

	// 调用已添加的hook函数
	ctx := context.Background()
	if hook, ok := hooks["category1"]["hookName1"]; ok {
		hook(ctx)
	}
}

// 定义hookFuncM类型
type hookFuncM map[string]map[string]func(ctx context.Context)

// newHookFuncM函数
func newHookFuncM() hookFuncM {
	return map[string]map[string]func(ctx context.Context){}
}

// 示例函数
func exampleHookFunc(ctx context.Context) {
	// 在此处理context和业务逻辑
	fmt.Println("Example hook function is called.")
}
