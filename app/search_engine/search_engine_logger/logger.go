package search_engine_logger

import (
	"github.com/sirupsen/logrus"
	log "github.com/wolanm/search-engine/logger"
	log2 "github.com/wolanm/search-engine/util/log"
)

var Logger *log.Logger

func InitLogger() {
	Logger, _ = log.NewLogger(log.LoggerConfig{
		Module: "search_engine_service",
		LogDir: log2.LogDir("search_engine_service"),
		Level:  logrus.InfoLevel,
	})
}
