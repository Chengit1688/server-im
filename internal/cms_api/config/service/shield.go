package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	adminModel "im/internal/cms_api/admin/model"
	adminRepo "im/internal/cms_api/admin/repo"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/common/constant"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
)

var ShieldService = new(shieldService)

type shieldService struct{}

func (s *shieldService) Add(ctx *gin.Context) {
	var (
		err       error
		req       model.ShieldReq
		data      model.ShieldWords
		adminUser *adminModel.Admin
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if err = util.CopyStructFields(&data, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" copy json, error: %v", err))
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	operationUserId := ctx.GetString("o_user_id")
	if operationUserId == "" {
		http.Failed(ctx, code.ErrUserPermissions)
		return
	}
	if adminUser, err = adminRepo.AdminRepo.GetByUserID(operationUserId); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " shield AdminRepo GetByUserID error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	data.OperationUser = adminUser.Username
	opt := configRepo.WhereOptionForShield{
		ShieldWords:  req.ShieldWords,
		DeleteStatus: constant.SwitchOff,
	}

	if total, _ := configRepo.ShieldRepo.Exists(opt); total != 0 {
		http.Failed(ctx, code.ErrShieldExist)
		return
	}
	if _, err = configRepo.ShieldRepo.Create(&data); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("Create, error: %v", err))
		http.Failed(ctx, code.ErrUnknown)
		return
	}

	mqtt.BroadcastMessage(req.OperationID, common.ConfigShieldPush, nil)
	http.Success(ctx)
}

func (s *shieldService) Update(ctx *gin.Context) {
	var (
		err  error
		req  model.ShieldReq
		data model.ShieldWords
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if req.ID == 0 {
		http.Failed(ctx, code.ErrBadRequest)
		return
	}
	if err = util.CopyStructFields(&data, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" Copy json, error: %v", err))
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	ep := configRepo.WhereOptionForShield{
		ShieldWords:  req.ShieldWords,
		DeleteStatus: constant.SwitchOff,
		Ext:          fmt.Sprintf("id != %d", req.ID),
	}
	if total, _ := configRepo.ShieldRepo.Exists(ep); total != 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("UpdateBy Exists, error: %v", err))
		http.Failed(ctx, code.ErrShieldExist)
		return
	}
	opt := configRepo.WhereOptionForShield{
		Id: req.ID,
	}
	if err = configRepo.ShieldRepo.UpdateById(opt, &data); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("UpdateBy, error: %v", err))
		http.Failed(ctx, err)
		return
	}

	mqtt.BroadcastMessage(req.OperationID, common.ConfigShieldPush, nil)
	http.Success(ctx)
}

func (s *shieldService) Delete(ctx *gin.Context) {
	var (
		err error
		req model.ShieldDeleteReq
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	opt := configRepo.WhereOptionForShield{
		Id: req.Id,
	}
	if err = configRepo.ShieldRepo.DeleteById(opt); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("DeleteById, error: %v", err))
		http.Failed(ctx, err)
		return
	}

	mqtt.BroadcastMessage(req.OperationID, common.ConfigShieldPush, nil)
	http.Success(ctx)
}

func (s *shieldService) GetList(ctx *gin.Context) {
	var (
		err   error
		count int64
		req   model.ShieldListReq
		resp  model.ShieldListResp
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(ctx, code.ErrBadRequest)
		return
	}
	if req.Status == 0 {
		req.Status = constant.SwitchOn
	}
	if req.DeleteStatus == 0 {
		req.DeleteStatus = constant.SwitchOff
	}
	resp.List, count, err = configRepo.ShieldRepo.GetList(req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("getlist, error: %v", err))
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(ctx, resp)
}
