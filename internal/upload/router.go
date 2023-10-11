package control

import (
	"github.com/gin-gonic/gin"
	minioService "im/internal/cms_api/third/service"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	baseRouter := router.Group("/api")
	thirdRouter := baseRouter.Group("/third")
	{
		thirdRouter.POST("/upload/v2", minioService.MinioService.UploadV2)
	}
	return router
}
