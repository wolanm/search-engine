package loading

import "github.com/wolanm/search-engine/config"

func Load() {
	// 配置加载
	config.InitConfig()

	// TODO: 日志初始化

	// TODO: 数据库连接初始化

	// TODO: redis 连接初始化

	// TODO: kafka 连接初始化
}
