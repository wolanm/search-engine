package indexplatform_logger

import (
	log "github.com/wolanm/search-engine/logger"
)

var Logger *log.Logger

func InitLogger() {
	Logger = log.SimpleLogger("index_platform_service")
}
