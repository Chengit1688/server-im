package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"im/internal/api/shopping/model"
	"im/internal/api/shopping/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var ShoppingService = new(shoppingService)

type shoppingService struct{}

func (s shoppingService) Search(ctx *gin.Context) {
	var (
		err   error
		req   model.SearchReq
		resp  model.SearchResp
		shops []model.SearchDTO
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	req.Check()
	if shops, resp.Count, err = repo.ShopRepo.SearchShop(req); err != nil {
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

func (s shoppingService) ApplyFor(ctx *gin.Context) {
	var (
		err error
		req model.ApplyForReq
	)

	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Infow(req.OperationID, "func", util.GetSelfFuncName(), "error", req)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := ctx.GetString("user_id")
	if _, err = repo.ShopRepo.AddShop(userID, req); err != nil {
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

func (s shoppingService) Update(ctx *gin.Context) {
	var (
		err error
		req model.ApplyForReq
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
	if _, err = repo.ShopRepo.UpdateShop(userID, req); err != nil {
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

func (s shoppingService) Detail(ctx *gin.Context) {
	var (
		err  error
		req  model.IDCommonReq
		resp model.ShopDetailResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if resp, err = repo.ShopRepo.FetchShop(req); err != nil {
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

func (s shoppingService) TeamMemberList(ctx *gin.Context) {
	var (
		err  error
		req  model.TeamListReq
		resp model.ShopTeamResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	req.Check()
	userID := ctx.GetString("user_id")
	if _, err = repo.ShopRepo.CheckShop(req.ShopID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
		return
	}
	if resp.List, resp.Count, err = repo.ShopRepo.FetchTeamList(userID, req); err != nil {
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

func (s shoppingService) JoinTeam(ctx *gin.Context) {
	var (
		err  error
		req  model.JoinTeamReq
		resp model.JoinTeamResp
		shop model.Shop
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
	if _, shop, err = repo.ShopRepo.JoinShopTeam(req); err != nil {
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

func (s shoppingService) RemoveTeam(ctx *gin.Context) {
	var (
		err error
		req model.RemoveTeamReq
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
	if err = repo.ShopRepo.RemoveShopTeam(req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrBadRequest, lang))
		return
	}
	http.Success(ctx)
}

func (s shoppingService) TeamMemberInfo(ctx *gin.Context) {
	var (
		err  error
		req  model.TeamMemberInfoReq
		resp model.TeamMemberInfoResp
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
	if resp, err = repo.ShopRepo.FetchShopTeamUser(req); err != nil {
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

func (s shoppingService) JoinTeamInfo(ctx *gin.Context) {
	var (
		err  error
		req  model.TeamLeaderInfoReq
		resp model.TeamLeaderInfoResp
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
	if resp, err = repo.ShopRepo.FetchShopTeamLeaderUser(req); err != nil {
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
