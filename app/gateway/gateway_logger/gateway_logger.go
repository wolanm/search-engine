package gateway_logger

import (
	"github.com/sirupsen/logrus"
	log "github.com/wolanm/search-engine/logger"
	log2 "github.com/wolanm/search-engine/util/log"
)

var Logger *log.Logger

func init() {
	Logger, _ = log.NewLogger(log.LoggerConfig{
		Module: "gateway",
		LogDir: log2.LogDir("gateway"),
		Level:  logrus.InfoLevel,
	})
}
