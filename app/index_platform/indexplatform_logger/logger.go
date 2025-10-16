package indexplatform_logger

import (
	"github.com/sirupsen/logrus"
	log "github.com/wolanm/search-engine/logger"
	log2 "github.com/wolanm/search-engine/util/log"
)

var Logger *log.Logger

func InitLogger() {
	Logger, _ = log.NewLogger(log.LoggerConfig{
		Module: "index_platform_service",
		LogDir: log2.LogDir("index_platform_service"),
		Level:  logrus.InfoLevel,
	})
}
