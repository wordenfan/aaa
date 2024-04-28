package api

import (
	"github.com/gin-gonic/gin"
)

type RouterInf interface {
	Route(r *gin.Engine)
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

type RegisterRouter struct{}

func (*RegisterRouter) RouteDistribute(ro RouterInf, r *gin.Engine) {
	ro.Route(r)
}

func InitRouter(r *gin.Engine) {
	rg := New()
	//rg.RouteDistribute(&user.RouterUser{}, r)
	//rg.RouteDistribute(&demo.RouterDemo{}, r)
	for _, router := range RouterSlice {
		rg.RouteDistribute(router, r)
	}
}
//============================================
var RouterSlice []RouterInf

func Register(ro ...RouterInf){
	RouterSlice = append(RouterSlice,ro...)
}