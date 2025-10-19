package search_engine_logger

import (
	"github.com/sirupsen/logrus"
	log "github.com/wolanm/search-engine/logger"
	log2 "github.com/wolanm/search-engine/util/log"
)

var Logger *log.Logger

func InitLogger() {
	Logger, _ = log.NewLogger(log.LoggerConfig{
		Module: "SearchEngineService",
		LogDir: log2.LogDir("SearchEngineService"),
		Level:  logrus.InfoLevel,
	})
}
