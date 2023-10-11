package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
)

var AppVersionService = new(appVersionService)

type appVersionService struct{}

func (s *appVersionService) Add(ctx *gin.Context) {
	var (
		err  error
		req  model.VersionReq
		data *model.AppVersion
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if err = util.Copy(req, &data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "Copy json error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	opt := configRepo.WhereOptionForVersion{
		Version:  req.Version,
		Platform: req.Platform,
	}
	if total, _ := configRepo.AppVersion.Exists(opt); total != 0 {
		http.Failed(ctx, code.ErrVersionExist)
		return
	}
	if _, err = configRepo.AppVersion.Create(data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "Create error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}

	http.Success(ctx)
}

func (s *appVersionService) Update(ctx *gin.Context) {
	var (
		err  error
		req  model.VersionReq
		data *model.AppVersion
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if err = util.Copy(req, &data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "Copy json error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	ep := configRepo.WhereOptionForVersion{
		Platform: req.Platform,
		Version:  req.Version,
		Ext:      fmt.Sprintf("id != %d", req.Id),
	}
	if total, _ := configRepo.AppVersion.Exists(ep); total != 0 {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "UpdateBy Exists error:", err)
		http.Failed(ctx, code.ErrVersionExist)
		return
	}
	opt := configRepo.WhereOptionForVersion{
		Id: req.Id,
	}
	if err = configRepo.AppVersion.UpdateById(opt, data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *appVersionService) UpdateStatus(ctx *gin.Context) {
	var (
		err error
		req model.VersionUpdateStatusReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(ctx, code.ErrBadRequest)
		return
	}
	if req.Status == constant.SwitchOn {
		opt := configRepo.WhereOptionForVersion{
			Status:   constant.SwitchOn,
			Platform: req.Platform,
		}
		if total, _ := configRepo.AppVersion.Exists(opt); total != 0 {
			http.Failed(ctx, code.ErrVersionRepeat)
			return
		}
	}
	opt := configRepo.WhereOptionForVersion{
		Id: req.Id,
	}
	if err = configRepo.AppVersion.UpdateById(opt, &model.AppVersion{Status: req.Status}); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *appVersionService) Delete(ctx *gin.Context) {
	var (
		err error
		req model.VersionDeleteReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(ctx, code.ErrBadRequest)
		return
	}

	opt := configRepo.WhereOptionForVersion{
		Id: req.Id,
	}
	if err = configRepo.AppVersion.UpdateById(opt, &model.AppVersion{DeleteStatus: constant.SwitchOn}); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *appVersionService) GetList(ctx *gin.Context) {
	var (
		err error
		req model.VersionListReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(ctx, code.ErrBadRequest)
		return
	}

	list, count, err := configRepo.AppVersion.GetList(req)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "GetList error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	ret := new(model.VersionListResp)
	_ = util.CopyStructFields(&ret.List, &list)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize

	http.Success(ctx, ret)
}
