package service

import (
	"github.com/gin-gonic/gin"
	"im/internal/api/user/model"
	"im/internal/api/user/repo"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"time"
)

var SignInV2Service = new(signInV2Service)

type signInV2Service struct{}

func (s *signInV2Service) SignIn(ctx *gin.Context) {
	var (
		err  error
		req  model.SignInReq
		sign model.SignInV2
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	UserID := ctx.GetString("user_id")
	ip := ctx.ClientIP()
	if sign, err = repo.SignInV2.FetchOne(UserID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	if sign.ID == 0 {
		if sign, err = repo.SignInV2.Add(UserID, ip); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrDB, lang))
			return
		}
	}
	currentTime := time.Now().Unix()
	daysDiff := util.GetDiffDaysBySecond(sign.LastTime, currentTime)
	if daysDiff <= 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrSignLog, lang))
		return
	}
	days := sign.ContinueDays + 1
	if daysDiff > 1 {
		days = 1
	}
	if sign, err = repo.SignInV2.SetContinueDays(UserID, days); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}

	http.Success(ctx)
}
