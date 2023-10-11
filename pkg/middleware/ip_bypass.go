package middleware

import (
	"bytes"
	"encoding/json"
	adminModel "im/internal/cms_api/admin/model"
	"im/internal/cms_api/ipwhitelist/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 配合后管IP白名单的中间件实现
func IPByPass() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("o_user_id")
		//处理登录接口 放行root
		for _, i := range WhiteListExclude {
			if util.KeyMatch(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				body, _ := ioutil.ReadAll(c.Request.Body)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
				var login adminModel.LoginReq
				json.Unmarshal(body, &login)
				if login.Username == "root" {
					c.Next()
					return
				}
			}
		}
		//处理非登录接口 放行root
		rootID := repo.IpWhiteListCache.GetRootID()
		if userID == rootID {
			c.Next()
			return
		} else {
			logger.Sugar.Errorw("not root", "userID:", userID, "rootID", rootID)
		}
		isOpen := repo.IpWhiteListCache.GetIsOpen()
		if isOpen == 1 {
			ip := c.ClientIP()
			status, err := repo.IpWhiteListCache.Exist(ip)
			if err != nil {
				logger.Sugar.Errorw(util.GetSelfFuncName(), "redis error:", err)
				c.Abort()
				http.Failed(c, code.ErrDB)
				return
			}
			if status {
				// ip 存在
				c.Next()
				return
			}
			c.Abort()
			http.Failed(c, code.ErrIPNotInWhiteList)
			return
		} else {
			c.Next()
		}
	}
}

// 白名单 特殊处理的路由列表
var WhiteListExclude = []UrlInfo{
	{Url: "/cms_api/admin/login", Method: "POST"},
}
