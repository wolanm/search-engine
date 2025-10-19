package main

import (
	"github.com/wolanm/search-engine/app/inverted_index/inverted_index_logger"
	"github.com/wolanm/search-engine/loading"
)

func main() {
	loading.Load()
	inverted_index_logger.InitLogger()

	// 启动消费者

}
