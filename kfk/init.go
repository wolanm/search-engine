package kfk

import (
	"github.com/IBM/sarama"
	"github.com/wolanm/search-engine/config"
	log "github.com/wolanm/search-engine/logger"
)

var GlobalKafkaClient sarama.Client

func InitKafka(logger *log.Logger) {
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaClient, err := sarama.NewClient(config.Conf.Kafka.Address, con)
	if err != nil {
		logger.Error(err)
		return
	}
	GlobalKafkaClient = kafkaClient

	logger.Info("create kafka client success on ", config.Conf.Kafka.Address)
}
