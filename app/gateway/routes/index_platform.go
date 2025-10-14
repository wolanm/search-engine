package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wolanm/search-engine/app/gateway/http"
)

func RegisterIndexPlatformRoutes(rg *gin.RouterGroup) {
	IndexPlatformGroup := rg.Group("/search_engine")
	{
		IndexPlatformGroup.POST("/build_index", http.BuildIndexByFiles)
		IndexPlatformGroup.POST("/upload_file", http.UploadFile)
	}
}
