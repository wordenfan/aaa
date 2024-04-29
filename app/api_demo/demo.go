package api_demo

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
	user_grpc "tk-boot-worden/app/user/cmd/rpc/pb"
	viper_conf "tk-boot-worden/config"
)

type HandlerDemo struct{}

func (*HandlerDemo) demoRequest(ctx *gin.Context) {
	//fmt.Println(" 我的测试 API! 222 ")
	//ctx.JSON(200, "我的测试 API! 222")
	postParam := "张三1"
	conn, err := grpc.Dial(viper_conf.GlobalConf.GrpcConf.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("gRPC连接失败: %v"+err.Error()) // 更改为Error级别，因为连接失败是一个重要问题
		return
	}
	defer conn.Close()
	UserClient := user_grpc.NewUcServiceClient(conn)
	c,cancel := context.WithTimeout(context.Background(),2*time.Second)
	defer cancel()
	grpcRsp,err := UserClient.UserInfo(c,&user_grpc.UserRequest{Name: postParam})
	if err != nil {
		zap.L().Error("client 连接失败: %v"+err.Error())
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK,fromError.Code())
	}
	ctx.JSON(http.StatusOK,grpcRsp.MyMessage)
}