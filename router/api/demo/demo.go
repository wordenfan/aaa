package demo

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type HandlerDemo struct{}

func init() {

}

func (*HandlerDemo) demoRequest(ctx *gin.Context) {
	fmt.Println(" 我的测试 API! ")
	ctx.JSON(200, "我的测试 API!")
}
