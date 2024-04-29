package demo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
	"tk-boot-worden/router/router"
)

var logger *rkentry.LoggerEntry

type RouterDemo struct{}

func (*RouterDemo) Route(r *gin.Engine) {
	h := &HandlerDemo{}
	user_group := r.Group("/v6", RouterMiddle())
	{
		user_group.GET("/demo_api", h.demoRequest)
	}
}

// ================================================
func init() {
	router.RouterSlice = append(router.RouterSlice, &RouterDemo{})
}

// ================================================
func RouterMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("路由分组中间件-before - 2")
		//logger.Info("路由分组中间件-before")
		// 可以在这里添加任何预处理逻辑，比如验证token、记录日志等
		// ...
		// 然后一定要调用c.Next()来传递给下一个处理器
		c.Next()
		fmt.Println("路由分组中间件-after - 2")
	}
}
