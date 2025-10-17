package main

import (
	"github.com/wolanm/search-engine/app/index_platform/indexplatform_logger"
	"github.com/wolanm/search-engine/app/index_platform/service"
	"github.com/wolanm/search-engine/config"
	"github.com/wolanm/search-engine/consts"
	pb "github.com/wolanm/search-engine/idl/pb/index_platform"
	"github.com/wolanm/search-engine/loading"
	"github.com/wolanm/search-engine/util/discovery"
	"google.golang.org/grpc"
	"net"
)

func runService() {
	grpcAddress := config.Conf.Services[consts.IndexPlatform].Addr[0]

	// 注册 etcd 服务
	node := &discovery.ServiceNode{
		ServiceName: config.Conf.Services[consts.IndexPlatform].Name,
		Endpoint:    grpcAddress,
	}
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister, err := discovery.NewServiceRegistry(etcdAddress, indexplatform_logger.Logger)
	if err != nil {
		panic(err)
	}
	err = etcdRegister.RegisterService(node, consts.LeaseTTL)
	if err != nil {
		panic(err)
	}

	// 注册 grpc 服务
	// TODO: 添加 otelgrpc 统计 grpc 调用数据, prometheus option
	server := grpc.NewServer()
	defer server.Stop()

	srvInstance := service.NewIndexPlatformSrv()
	pb.RegisterIndexPlatformServiceServer(server, srvInstance)

	// TODO: 注册 prometheus 服务
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	indexplatform_logger.Logger.Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}

func main() {
	// 配置加载
	loading.Load()

	// 分词器初始化
	anylzer.InitSeg()

	// 日志初始化
	indexplatform_logger.InitLogger()

	// TODO: 注册 tracer

	// 启动 grpc 服务
	runService()
}
