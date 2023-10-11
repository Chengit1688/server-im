package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"im/internal/api/operator/model"
	"im/internal/api/operator/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var OperatorService = new(operatorService)

type operatorService struct{}

func (s operatorService) Search(ctx *gin.Context) {
	var (
		err   error
		req   model.SearchOperatorReq
		resp  model.SearchOperatorResp
		shops []model.SearchDTO
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	req.Check()
	if shops, resp.Count, err = repo.OperatorRepo.SearchShop(req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	for _, shop := range shops {
		resp.List = append(resp.List, model.SearchInfo{
			ShopID:          shop.ID,
			ShopName:        shop.Name,
			DecorationScore: shop.DecorationScore,
			QualityScore:    shop.QualityScore,
			ServiceScore:    shop.ServiceScore,
			ShopStar:        shop.Star,
			ShopLocation:    shop.Address,
			ShopType:        shop.ShopType,
			ShopDistance:    shop.Distance,
			ShopIcon:        strings.Split(shop.Image, ","),
			Latitude:        shop.Latitude,
			Longitude:       shop.Longitude,
		})
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(ctx, resp)
}

func (s operatorService) ApplyFor(ctx *gin.Context) {
	var (
		err error
		req model.OperatorApplyForReq
	)

	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Infow(req.OperationID, "func", util.GetSelfFuncName(), "error", req)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := ctx.GetString("user_id")
	if _, err = repo.OperatorRepo.AddShop(userID, req); err != nil {
		if err == code.ErrShopExists {
			logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrShopExists, lang))
			return
		}
		if err == code.ErrInviteCodeExists {
			logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrInviteCodeExists, lang))
			return
		}
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	http.Success(ctx)
}

func (s operatorService) Update(ctx *gin.Context) {
	var (
		err error
		req model.OperatorApplyForReq
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.ShopID == 0 {
		logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := ctx.GetString("user_id")
	if _, err = repo.OperatorRepo.UpdateShop(userID, req); err != nil {
		if err == code.ErrShopNotExists {
			logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
			return
		}
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}

	http.Success(ctx)
}

func (s operatorService) Detail(ctx *gin.Context) {
	var (
		err  error
		req  model.OperatorIDCommonReq
		resp model.OperatorDetailResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if resp, err = repo.OperatorRepo.FetchShop(req); err != nil {
		if err == code.ErrShopNotExists {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
			return
		}
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	http.Success(ctx, resp)
}

func (s operatorService) TeamMemberList(ctx *gin.Context) {
	var (
		err  error
		req  model.OperatorTeamListReq
		resp model.OperatorTeamResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	req.Check()
	userID := ctx.GetString("user_id")
	if _, err = repo.OperatorRepo.CheckShop(req.ShopID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
		return
	}
	if resp.List, resp.Count, err = repo.OperatorRepo.FetchTeamList(userID, req); err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrDB, lang))
			return
		}
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(ctx, resp)
}

func (s operatorService) JoinTeam(ctx *gin.Context) {
	var (
		err  error
		req  model.OperatorJoinTeamReq
		resp model.OperatorJoinTeamResp
		shop model.Operator
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.InviteCode == "" && req.ShopID == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.UserID == "" {
		req.UserID = ctx.GetString("user_id")
	}
	if _, shop, err = repo.OperatorRepo.JoinShopTeam(req); err != nil {
		if err == code.ErrShopNotExists {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
			return
		}
		if err == code.ErrInviteCodeNotExists {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrInviteCodeNotExists, lang))
			return
		}
		if err == code.ErrShopTeamMemberExists {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrShopTeamMemberExists, lang))
			return
		}
		if err == code.ErrShopExSelf {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrShopExSelf, lang))
			return
		}
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}
	resp.ShopID = shop.ID
	resp.TeamName = shop.Name

	http.Success(ctx, resp)
}

func (s operatorService) RemoveTeam(ctx *gin.Context) {
	var (
		err error
		req model.OperatorRemoveTeamReq
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.UserID == "" {
		req.UserID = ctx.GetString("user_id")
	}
	if err = repo.OperatorRepo.RemoveShopTeam(req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	http.Success(ctx)
}

func (s operatorService) TeamMemberInfo(ctx *gin.Context) {
	var (
		err  error
		req  model.OperatorTeamMemberInfoReq
		resp model.OperatorTeamMemberInfoResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.UserID == "" {
		req.UserID = ctx.GetString("user_id")
	}
	if resp, err = repo.OperatorRepo.FetchShopTeamUser(req); err != nil {
		if err == code.ErrUserNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrUserNotFound, lang))
			return
		}
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrDB, lang))
		return
	}

	http.Success(ctx, resp)
}

func (s operatorService) JoinTeamInfo(ctx *gin.Context) {
	var (
		err  error
		req  model.OperatorTeamLeaderInfoReq
		resp model.OperatorTeamLeaderInfoResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.UserID == "" {
		req.UserID = ctx.GetString("user_id")
	}
	if resp, err = repo.OperatorRepo.FetchShopTeamLeaderUser(req); err != nil {
		if err == code.ErrUserNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrUserNotFound, lang))
			return
		}
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(ctx, response.GetError(response.ErrDB, lang))
			return
		}
	}

	http.Success(ctx, resp)
}
