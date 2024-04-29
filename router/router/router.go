package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/resolver"
	"log"
	viper_conf "tk-boot-worden/config"
	"tk-boot-worden/discovery"
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

var RouterSlice []RouterInf

func InitRouter(r *gin.Engine) {
	rg := New()
	//rg.RouteDistribute(&user.RouterUser{}, r)
	//rg.RouteDistribute(&demo.RouterDemo{}, r)
	for _, router := range RouterSlice {
		rg.RouteDistribute(router, r)
	}
}

// ===================================
func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(viper_conf.GlobalConf.EtcdConf.AddrS, logs.LG)
	resolver.Register(etcdRegister)
	info := discovery.Server{
		Name:    viper_conf.GlobalConf.GrpcConf.Name,
		Addr:    viper_conf.GlobalConf.GrpcConf.Addr,
		Version: viper_conf.GlobalConf.GrpcConf.Version,
		Weight:  viper_conf.GlobalConf.GrpcConf.Weight,
	}
	r := discovery.NewRegister(viper_conf.GlobalConf.EtcdConf.AddrS, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
