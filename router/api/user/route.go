package user

import (
	"github.com/gin-gonic/gin"
)

type RouterUser struct{}

func (*RouterUser) Route(r *gin.Engine) {
	h := &HandlerUser{}
	r.GET("/v1/captcha", h.getcaptcha)
	r.GET("/v1/get_redis", h.GetRedis)
	r.GET("/v1/set_redis", h.SetRedis)
}
