package tools

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	order_grpc "tk-boot-worden/app/order/cmd/rpc/pb"
	user_grpc "tk-boot-worden/app/user/cmd/rpc/pb"
	viper_conf "tk-boot-worden/config"
)

type grpcConfig struct {
	Addr         string
	RegisterFunc func(server *grpc.Server)
	Version      string
	Weight       int64
}

type UserServer struct {
	user_grpc.UnimplementedUcServiceServer
}
type OrderServer struct {
	order_grpc.UnimplementedOrderServiceServer
}

func (server *UserServer) UserInfo(_ context.Context, req *user_grpc.UserRequest) (*user_grpc.UserResponse, error) {
	assumeError := true
	if assumeError {
		//return nil, status.Error(1001,"假定的错误执行了")
		return nil, transferToGrpcError(IlLegalUser)
	}
	return &user_grpc.UserResponse{
		MyMessage: "hello! " + req.Name,
	}, nil
}

func (server *OrderServer) OrderInfo(_ context.Context, req *order_grpc.OrderRequest) (*order_grpc.OrderResponse, error) {
	p := &order_grpc.OrderResponse{
		Id:    req.GetId(),
		Name:  "worden",
		Email: "rs@example.com",
		Phones: []*order_grpc.OrderResponse_PhoneNumber{
			{Number: "555-4321", Type: order_grpc.OrderResponse_HOME},
		},
	}
	return p, nil
}
func InitGrpcServer() *grpc.Server {
	c := grpcConfig{
		Addr: viper_conf.GlobalConf.GrpcConf.Addr,
		RegisterFunc: func(server *grpc.Server) {
			user_grpc.RegisterUcServiceServer(server, &UserServer{})
			order_grpc.RegisterOrderServiceServer(server, &OrderServer{})
		},
		Version: viper_conf.GlobalConf.GrpcConf.Version,
		Weight:  viper_conf.GlobalConf.GrpcConf.Weight,
	}
	srv := grpc.NewServer()
	c.RegisterFunc(srv)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		zap.L().Debug("cannot listen")
	}
	go func() {
		err = srv.Serve(lis)
		if err != nil {
			zap.L().Debug("grpc server started error: " + err.Error())
			return
		}
	}()
	return srv
}
