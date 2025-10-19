package kfk

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"github.com/wolanm/search-engine/app/inverted_index/inverted_index_logger"
	"github.com/wolanm/search-engine/app/inverted_index/service"
	"github.com/wolanm/search-engine/config"
	"github.com/wolanm/search-engine/kfk"
	"github.com/wolanm/search-engine/types"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type InvertedIndexConsumer struct {
	Ready chan bool
}

func (consumer *InvertedIndexConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *InvertedIndexConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *InvertedIndexConsumer) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				inverted_index_logger.Logger.Info("message channel closed")
				return nil
			}

			fileInfo := new(types.FileInfo)
			_ = fileInfo.UnmarshalJSON(msg.Value)

			// 调用 mapreduce 处理
			go service.BuildInvertedIndex(fileInfo)

			inverted_index_logger.Logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s",
				string(msg.Value), msg.Timestamp, msg.Topic)
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			inverted_index_logger.Logger.Info("session context done")
			return nil
		}
	}

	return nil
}

// InvertedIndexKafkaConsume 倒排索引消费者建立
func InvertedIndexKafkaConsume(ctx context.Context, topic, group, assignor string) (err error) {
	keepRunning := true
	inverted_index_logger.Logger.Info("start a new sarama consumer")
	sarama.Logger = inverted_index_logger.Logger

	// 创建消费者组 TODO:多文件上传时，可以添加多个消费者提升性能
	consumer := InvertedIndexConsumer{
		Ready: make(chan bool),
	}

	consumerConfig := kfk.GetDefaultConsumeConfig(assignor) // 用来停止 Consume
	cancelCtx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup(config.Conf.Kafka.Address, group, consumerConfig)
	if err != nil {
		inverted_index_logger.Logger.Error("create consumer group failed: ", err)
		cancel()
		return err
	}

	//isConsumptionPause := false
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// 持续消费
	go func() {
		defer wg.Done()
		for {
			if err = client.Consume(cancelCtx, []string{topic}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				inverted_index_logger.Logger.Errorf("Error from consumer: %v", err)
			}

			// 如果是因为调用 cancel() 关闭则直接退出
			if cancelCtx.Err() != nil {
				return
			}

			// 重新创建chan 以便消费重试
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready
	inverted_index_logger.Logger.Info("inverted index consumer start running")

	// SIGUSR1 是 linux 系统的信号
	//sigusr1 := make(chan os.Signal, 1) // 用来控制消费的暂停和恢复
	//signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)

	for keepRunning {
		select {
		case <-cancelCtx.Done():
			inverted_index_logger.Logger.Infof("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			inverted_index_logger.Logger.Infof("terminating: via term signal")
			keepRunning = false
			//case <-sigusr1:
			//	isConsumptionPause = !isConsumptionPause
			//	toggleConsumptionFlow(client, isConsumptionPause)
		}
	}

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		inverted_index_logger.Logger.Error("close consumer group failed: ", err)
		return
	}

	return
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused bool) {
	if isPaused {
		client.PauseAll()
		inverted_index_logger.Logger.Info("pause consumption")
	} else {
		client.ResumeAll()
		inverted_index_logger.Logger.Info("resume consumption")
	}
}
