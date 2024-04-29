package tools

import "fmt"

//全局配置变量

type GinResponse struct {
	Code    uint32      `json:"code"`
	Msg     string      `json:"msg"`
	Details []interface{} `json:"details"`
}

func (m GinResponse) PrintResponse() string {
	return fmt.Sprintf("%d-%s", m.Code, m.Msg)
}