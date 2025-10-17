package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wolanm/search-engine/app/gateway/common"
	"github.com/wolanm/search-engine/app/gateway/gateway_logger"
	"github.com/wolanm/search-engine/app/gateway/response"
	"github.com/wolanm/search-engine/app/gateway/rpc"
	"net/http"
)

func BuildIndexByFiles(ctx *gin.Context) {
	response.HTTPResponse(ctx, http.StatusOK, 0, "BuildIndexByFiles called")
}

func UploadFile(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if nil == fileHeader {
		gateway_logger.Logger.Error("UploadFile get file failed: ", err)
		response.HTTPResponse(ctx, http.StatusBadRequest, 0, "upload file error")
		return
	}

	rpcResp, err := rpc.UploadFile(ctx, file, common.FileInfo{
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
	})
	if err != nil {
		gateway_logger.Logger.Error("rpc failed: ", err)
		response.HTTPResponse(ctx, http.StatusInternalServerError, 0, "server error")
		return
	}

	response.HTTPResponse(ctx, http.StatusOK, int(rpcResp.Code), rpcResp.Message)
}
