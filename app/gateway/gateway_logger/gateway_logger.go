package gateway_logger

import log "github.com/wolanm/search-engine/logger"

var Logger *log.Logger

func init() {
	Logger = log.SimpleLogger("gateway")
}
