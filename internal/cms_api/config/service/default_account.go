package service

import (
	"github.com/gin-gonic/gin"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
)

var DefaultAccountService = new(defaultAccountService)

type defaultAccountService struct{}

func (s *defaultAccountService) AddOrUpdateFriend(ctx *gin.Context) {
	var (
		err  error
		t    int64
		req  model.DefaultFriendReq
		user *userModel.User
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	operationUserId := ctx.GetString("o_user_id")
	if operationUserId == "" {
		http.Failed(ctx, code.ErrUserPermissions)
		return
	}
	if req.Id != 0 {
		if t, _ = configRepo.DefaultFriendRepo.Exists(configRepo.WhereOptionForDefaultFriend{Id: req.Id}); t == 0 {
			http.Failed(ctx, code.ErrDefaultFriendNotFound)
			return
		}
	} else {
		if t, _ = configRepo.DefaultFriendRepo.Exists(configRepo.WhereOptionForDefaultFriend{UserId: req.UserId}); t > 0 {
			http.Failed(ctx, code.ErrDefaultFriendFound)
			return
		}
	}
	userOp := userRepo.WhereOption{
		UserId: req.UserId,
		Status: 1,
	}
	if user, err = userRepo.UserRepo.GetByUserID(userOp); err != nil {
		http.Failed(ctx, code.ErrUserIdNotExist)
		return
	}
	data := model.DefaultFriend{
		UserPrimaryId:   user.ID,
		UserId:          user.UserID,
		GreetMsg:        req.GreetMsg,
		Remarks:         req.Remarks,
		OperationUserId: operationUserId,
	}
	if _, err = configRepo.DefaultFriendRepo.CreateOrUpdate(&data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend configRepo GetByUserID error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}

	http.Success(ctx)
}

func (s *defaultAccountService) DeleteFriend(ctx *gin.Context) {
	var (
		err error
		req model.DefaultFriendDeleteReq
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if req.UserId != "" {
		if err = configRepo.DefaultFriendRepo.DeleteByUserId(req.UserId); err != nil {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend DeleteByUserId error:", err)
			http.Failed(ctx, err)
			return
		}
	}
	if err = configRepo.DefaultFriendRepo.DeleteById(req.Id); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend DeleteById error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *defaultAccountService) GetFriendList(ctx *gin.Context) {
	var (
		err   error
		count int64
		req   model.DefaultFriendListReq
		resp  model.DefaultFriendListResp
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	resp.List, count, err = configRepo.DefaultFriendRepo.GetList(req)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " default_friend GetList error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(ctx, resp)
}
