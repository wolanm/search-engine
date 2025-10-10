package main

import (
	"fmt"
	"github.com/wolanm/search-engine/app/gateway/routes"
	"github.com/wolanm/search-engine/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startListen() {
	ginRouter := routes.NewRouter()
	server := &http.Server{
		Addr:           config.Conf.Server.Port,
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("绑定HTTP到 %s 失败！可能是端口已经被占用，或用户权限不足 \n", config.Conf.Server.Port)
		fmt.Println(err)
		return
	}
}

func main() {
	// 配置加载
	config.InitConfig()

	// rpc 初始化

	// etcd 注册

	// 启动 web 服务
	go startListen()

	fmt.Printf("gateway listen on :%v \n", config.Conf.Server.Port)
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		s := <-osSignals
		fmt.Println("exit! ", s)
	}
}
