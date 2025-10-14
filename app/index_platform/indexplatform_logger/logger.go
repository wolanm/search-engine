package indexplatform_logger

import (
	log "github.com/wolanm/search-engine/logger"
)

var Logger *log.Logger

func initLogger() {
	Logger = log.SimpleLogger("index_platform_service")
}
