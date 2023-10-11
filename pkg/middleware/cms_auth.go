package middleware

import (
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
	"github.com/gin-gonic/gin"
)

// Auth jwt验证
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/cms_api/swagger") {
			return
		}
		var jwtAuthExclude bool
		for _, i := range JwtAuthExclude {
			if util.KeyMatch(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				jwtAuthExclude = true
				break
			}
		}
		if jwtAuthExclude {
			// 排除 放行
			logger.Sugar.Info("jwt allow ", util.GetSelfFuncName(), c.Request.Method, c.Request.URL.Path)
			c.Next()
			return
		}
		userID, roleKey, _, err := util.CmsParseToken(c.Request.Header.Get("token"))
		if err != nil {
			logger.Sugar.Error(c.Request.Header.Get("token"), util.GetSelfFuncName(), "jwt parse token error:", err)
			c.Abort()
			http.Failed(c, code.ErrUnauthorized)
			return
		}
		c.Set("o_user_id", userID)
		c.Set("o_role_key", roleKey)
		c.Next()
	}
}

// AuthCheckRole 权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		userID, roleKey, _, _ := util.CmsParseToken(c.Request.Header.Get("token"))

		e := util.GetEnforcer()
		var res, casbinExclude bool

		//检查权限
		if roleKey == "admin" {
			res = true
			// admin 放行
			c.Next()
			return
		}
		for _, i := range CasbinExclude {
			if util.KeyMatch(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				casbinExclude = true
				break
			}
		}
		if casbinExclude {
			// 排除 放行
			logger.Sugar.Info("Casbin exclusion, no validation method ", util.GetSelfFuncName(), " ", c.Request.Method, " ", c.Request.URL.Path)
			c.Next()
			return
		}
		res, err = e.Enforce(roleKey, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			logger.Sugar.Error("AuthCheckRole error: ", err, " method: ", c.Request.Method, " path: ", c.Request.URL.Path, " user id: ", userID)
			http.Failed(c, code.ErrNoPermission)
			return
		}

		if res {
			logger.Sugar.Info("Casbin Permit Access ", util.GetSelfFuncName(), " role ", roleKey, " method: ", c.Request.Method, " path: ", c.Request.URL.Path)
			c.Next()
		} else {
			logger.Sugar.Error("Casbin Deny Access: ", util.GetSelfFuncName(), " role ", roleKey, " method: ", c.Request.Method, " path:", c.Request.URL.Path, " user id: ", userID)
			http.Failed(c, code.ErrNoPermission)
			c.Abort()
			return
		}

	}
}

type UrlInfo struct {
	Url    string
	Method string
}

// CasbinExclude casbin 排除的路由列表
var CasbinExclude = []UrlInfo{
	{Url: "/cms_api/admin/login", Method: "POST"},
	{Url: "/cms_api/admin/logout", Method: "POST"},
}

// Jwt 认证 排除的路由列表
var JwtAuthExclude = []UrlInfo{
	{Url: "/cms_api/admin/login", Method: "POST"},
	{Url: "/cms_api/config/info", Method: "GET"},
}
