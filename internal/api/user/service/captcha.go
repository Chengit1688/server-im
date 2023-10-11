package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	settingModel "im/internal/api/setting/model"
	"im/internal/api/user/usecase"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
)

var CaptchaService = new(captchaService)

type captchaService struct{}

func (s *captchaService) GetCaptcha(ctx *gin.Context) {
	var (
		err error
		req settingModel.CaptchaReq
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" bind json , error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	data, err := usecase.GetCaptchaFactory().GetService(req.CaptchaType).Get()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" GetCaptcha copy , error: %v", err))
		http.Failed(ctx, err)
		return
	}
	http.Success(ctx, data)
}

func (s *captchaService) CheckCaptcha(ctx *gin.Context) {
	var (
		err error
		req settingModel.CheckCaptchaReq
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" bind json , error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	ser := usecase.GetCaptchaFactory().GetService(req.CaptchaType)
	if err = ser.Check(req.Token, req.PointJson); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" CheckCaptcha copy , error: %v", err))
		http.Failed(ctx, err)
		return
	}
	http.Success(ctx)
}
