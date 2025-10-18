package main

import (
	"github.com/wolanm/search-engine/app/search_engine/search_engine_logger"
	"github.com/wolanm/search-engine/config"
)

func main() {
	config.InitConfig()
	search_engine_logger.InitLogger()

	// TODO: 服务注册与发现

	// TODO: grpc client 初始化

	// TODO: grpc server 初始化
}
