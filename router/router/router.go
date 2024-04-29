package router

import (
	"github.com/gin-gonic/gin"
	viper_conf "tk-boot-worden/config"
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
		Name:    viper_conf.GrpcConfig.NameName,
		Addr:    viper_conf.GrpcConfig.Addr,
		Version: viper_conf.GrpcConfig.Version,
		Weight:  viper_conf.Grpcconfig.Weight,
	}
	r := discovery.NewRegister(config.AppConf.Etcdconfig.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
