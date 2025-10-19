package indexplatform_logger

import (
	"github.com/sirupsen/logrus"
	"github.com/wolanm/search-engine/consts"
	log "github.com/wolanm/search-engine/logger"
	log2 "github.com/wolanm/search-engine/util/log"
)

var Logger *log.Logger

func InitLogger() {
	Logger, _ = log.NewLogger(log.LoggerConfig{
		Module: consts.IndexPlatformService,
		LogDir: log2.LogDir(consts.IndexPlatformService),
		Level:  logrus.InfoLevel,
	})
}
