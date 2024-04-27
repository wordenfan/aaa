package demo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rkentry "github.com/rookie-ninja/rk-entry/v2/entry"
)

var logger *rkentry.LoggerEntry

type RouterDemo struct{}

//redisGroup.GET("/demo_api", demoRequest)

func (*RouterDemo) Route(r *gin.Engine) {
	h := &HandlerUser{}
	user_group := r.Group("/v6")
	{
		user_group.GET("/demo_api", h.demoRequest)
		user_group.Use(RouterMiddle())
	}

	//r.GET("/v3/demo_api", h.demoRequest)
}

// ================================================
func RouterMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("路由分组中间件-before")
		logger.Info("路由分组中间件-before")
		// 可以在这里添加任何预处理逻辑，比如验证token、记录日志等
		// ...
		// 然后一定要调用c.Next()来传递给下一个处理器
		c.Next()
		fmt.Println("路由分组中间件-after")
	}
}
