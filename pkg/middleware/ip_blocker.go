package middleware

import (
	"im/internal/cms_api/ipblacklist/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

// 配合IP黑名单的中间件实现
func IPBlocker() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		route := c.Request.URL.Path
		if route == "/api/setting/config" {
			return
		}
		status, err := repo.IpBlackListCache.Exist(ip)
		if err != nil {
			logger.Sugar.Errorw(util.GetSelfFuncName(), "redis error:", err)
			c.Abort()
			http.Failed(c, code.ErrDB)
			return
		}
		if status {
			// ip 存在
			c.Abort()
			http.Failed(c, code.ErrIpBlocked)
			return
		}
		c.Next()
	}
}
