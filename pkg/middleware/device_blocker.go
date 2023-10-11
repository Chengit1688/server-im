package middleware

import (
	"im/internal/cms_api/user/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

// 配合封禁管理 设备封禁的功能实现
func DeviceBlocker() gin.HandlerFunc {
	return func(c *gin.Context) {
		device_id := c.GetHeader("device_id")
		//logger.Sugar.Debugw(util.GetSelfFuncName(), "device_id:", device_id)
		route := c.Request.URL.Path
		if route == "/api/setting/config" {
			return
		}
		if len(device_id) == 0 {
			c.Next()
		}
		status, err := repo.DeviceListCache.Exist(device_id)
		if err != nil {
			logger.Sugar.Errorw(util.GetSelfFuncName(), "redis error:", err)
			c.Abort()
			http.Failed(c, code.ErrDB)
			return
		}
		if status {
			// device_id 存在
			c.Abort()
			http.Failed(c, code.ErrDeviceIDBlocked)
			return
		}
		c.Next()
	}
}

//防止重复请求
//func RepeatSub()gin.HandlerFunc{
//	return func(c *gin.Context) {
//		c.GetQuery("operation_id")
//		c.GetPostForm()
//	}
//}