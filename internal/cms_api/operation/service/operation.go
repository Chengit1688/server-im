package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiUserModel "im/internal/api/user/model"
	"im/internal/cms_api/operation/model"
	"im/internal/cms_api/operation/repo"
	userRepo "im/internal/cms_api/user/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
)

var OperationService = new(operationService)

type operationService struct{}

func (s *operationService) RegistrationStatistics(ctx *gin.Context) {
	var (
		err  error
		req  model.RegistrationStatisticsReq
		resp model.RegistrationStatisticsResp
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " RegistrationStatistics bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if resp.List, err = repo.OperationRepo.GetRegistrationStatistics(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " RegistrationStatistics GetRegistrationStatistics error:", err)
		http.Failed(ctx, code.ErrBadSignLog)
		return
	}

	http.Success(ctx, resp)
}

func (s *operationService) InviteCodeStatistics(ctx *gin.Context) {
	var (
		err   error
		count int64
		req   model.InviteCodeStatisticsReq
		resp  model.InviteCodeStatisticsResp
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " InviteCodeStatistics bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if resp.List, count, err = repo.OperationRepo.GetInviteCodeStatistics(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " InviteCodeStatistics error:", err)
		http.Failed(ctx, code.ErrBadSignLog)
		return
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(ctx, resp)
}

func (s *operationService) InviteCodeStatisticsDetails(ctx *gin.Context) {
	var (
		err   error
		count int64
		req   model.InviteCodeStatisticsDetailsReq
		resp  model.InviteCodeStatisticsDetailsResp
	)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if resp.List, count, err = repo.OperationRepo.GetInviteCodeStatisticsDetails(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " GetInviteCodeStatisticsDetails error:", err)
		http.Failed(ctx, code.ErrBadSignLog)
		return
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(ctx, resp)
}

func (s *operationService) GetSuggestionList(ctx *gin.Context) {
	var (
		err  error
		req  model.SuggestionReq
		resp model.SuggestionInfoResp
	)
	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " suggestion bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	req.Pagination.Check()
	opt := repo.WhereOptionForSuggestion{
		Account:    req.Account,
		NickName:   req.NickName,
		UserId:     req.UserId,
		Content:    req.Content,
		Brand:      req.Brand,
		Platform:   req.Platform,
		BeginDate:  req.BeginDate,
		EndDate:    req.EndDate,
		Pagination: req.Pagination,
	}
	list, count, err := repo.SuggestionRepo.GetList(opt)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " iGetSuggestionList error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	_ = util.CopyStructFields(&resp.List, &list)
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(ctx, resp)
}

func (s *operationService) OnlineUsers(ctx *gin.Context) {
	var (
		req  model.OnlineUsersReq
		resp model.OnlineUsersResp
		err  error
	)

	if err = ctx.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, code.ErrBadRequest)
		return
	}

	req.Check()
	resp.Pagination = req.Pagination

	var users []apiUserModel.User

	if users, resp.Count, err = userRepo.UserRepo.GetOnlineUsers(req.Keyword, req.Offset, req.Limit); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get online users error, error: %v", err))
		http.Failed(ctx, code.ErrDB)
		return
	}

	resp.TotalCount = resp.Count

	if req.Keyword != "" {
		if resp.TotalCount, err = userRepo.UserRepo.GetOnlineUsersCount(); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get online users count error, error: %v", err))
			http.Failed(ctx, code.ErrDB)
			return
		}
	}

	if len(users) != 0 {
		for _, user := range users {
			var info model.OnlineUserInfo
			info.UserID = user.UserID
			info.Nickname = user.NickName
			info.LastLoginTime = user.LatestLoginTime
			resp.List = append(resp.List, info)
		}
	}

	http.Success(ctx, resp)
}
