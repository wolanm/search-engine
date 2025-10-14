package response

import "github.com/gin-gonic/gin"

func HTTPResponse(ctx *gin.Context, httpCode int, code int, message string) {
	ctx.JSON(httpCode, gin.H{
		"code":    code,
		"message": message,
	})
}

func HTTPResponseWithData(ctx *gin.Context, httpCode int, code int, message string, data any) {
	ctx.JSON(httpCode, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
