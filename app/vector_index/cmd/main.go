package main

import (
	"github.com/wolanm/search-engine/app/vector_index/vector_index_logger"
	"github.com/wolanm/search-engine/loading"
)

func main() {
	loading.Load()

	vector_index_logger.InitLogger()

	// TODO: kafka 初始化

}
