package main

import (
	"context"
	"github.com/wolanm/search-engine/app/inverted_index/analyzer"
	"github.com/wolanm/search-engine/app/inverted_index/inverted_index_logger"
	"github.com/wolanm/search-engine/app/inverted_index/service"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/kfk/index_consumer"
	"github.com/wolanm/search-engine/loading"
	"sync"
)

func main() {
	loading.Load()

	inverted_index_logger.InitLogger()

	// 分词器初始化
	analyzer.InitSeg()

	// 启动消费者

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := index_consumer.IndexKafkaConsume(context.Background(), consts.KafkaIndexTopic, consts.InvertedIndexGroupID,
			consts.KafkaAssignorRoundRobin, service.BuildInvertedIndex)
		if err != nil {
			inverted_index_logger.Logger.Error("stop inverted index consume: ", err)
		}
	}()

	wg.Wait()

	// TODO: 注册倒排索引 rpc 服务
}
