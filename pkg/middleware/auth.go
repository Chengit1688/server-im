package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	userRepo "im/internal/api/user/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
)

func OAuth() gin.HandlerFunc {
	whiteList := []string{
		"/api/user/login",
		"/api/user/register",
		"/api/setting/config",
		"/api/setting/version",
		"/api/setting/sms",
		"/api/captcha/get",
		"/api/captcha/check",
		"/api/user/forgot_password",
		"/api/chat/push_message",
		"/api/user/verify_code",
		"/api/swagger/",
		"/api/setting/domain_list",
		"/api/setting/privacy_policy",
		"/api/setting/user_agreement",
	}
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token") //token
		route := ctx.Request.URL.Path
		//匹配白名单
		rIndex := util.IndexOf(whiteList, route)
		if rIndex != -1 || strings.Contains(route, "/swagger/") {
			return
		}
		if token == "" {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", "token is empty!")
			ctx.Abort()
			http.Failed(ctx, code.ErrUnauthorized)
			return
		}
		userId, err := userRepo.UserCache.TokenInfo(token)
		if err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("parse token error: %v", err))
			ctx.Abort()
			http.Failed(ctx, code.ErrUnauthorized)
			return
		}
		if userId == "" {
			ctx.Abort()
			http.Failed(ctx, code.ErrUnauthorized)
			return
		}
		ctx.Set("user_id", userId)
	}
}
