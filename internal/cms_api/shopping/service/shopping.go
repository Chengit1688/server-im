package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	friendUseCase "im/internal/api/friend/usecase"
	shoppingModel "im/internal/api/shopping/model"
	shoppingRepo "im/internal/api/shopping/repo"
	configModel "im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/internal/cms_api/shopping/model"
	"im/internal/cms_api/shopping/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var ShoppingService = new(shoppingService)

type shoppingService struct{}

func (s *shoppingService) ShopList(c *gin.Context) {
	var (
		err   error
		count int64
		req   model.ShopListReq
		resp  model.ShopListResp
		shops []shoppingModel.Shop
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	req.Check()
	if shops, count, err = repo.ShopRepo.FetchList(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	for _, shop := range shops {
		resp.List = append(resp.List, model.ShopListInfo{UserInfo: model.UserInfo{
			UserID:      shop.CreatorUser.UserID,
			Account:     shop.CreatorUser.Account,
			PhoneNumber: shop.CreatorUser.PhoneNumber,
			CountryCode: shop.CreatorUser.CountryCode,
			FaceURL:     shop.CreatorUser.FaceURL,
			BigFaceURL:  shop.CreatorUser.BigFaceURL,
			Gender:      shop.CreatorUser.Gender,
			NickName:    shop.CreatorUser.NickName,
		}, ShopInfo: model.ShopInfo{
			ShopID:          shop.ID,
			ShopName:        shop.Name,
			ShopLocation:    shop.Longitude,
			ShopType:        shop.ShopType,
			DecorationScore: shop.DecorationScore,
			Star:            shop.Star,
			QualityScore:    shop.QualityScore,
			ServiceScore:    shop.ServiceScore,
			ShopIcon:        strings.Split(shop.Image, ","),
			Status:          shop.Status,
			Longitude:       shop.Longitude,
			Latitude:        shop.Latitude,
			Address:         shop.Address,
			License:         shop.License,
			CreatedAt:       shop.CreatedAt,
			Description:     shop.Description,
		}})
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(c, resp)
}

func (s *shoppingService) MemberList(ctx *gin.Context) {
	var (
		err  error
		req  model.MemberListReq
		resp model.MemberListResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	req.Check()
	if _, err = shoppingRepo.ShopRepo.CheckShop(req.ShopID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
		return
	}
	if resp.List, resp.Count, err = repo.ShopRepo.FetchTeamList(req); err != nil {
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

func (s *shoppingService) Approve(c *gin.Context) {
	var (
		err  error
		req  model.ApproveReq
		shop shoppingModel.Shop
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	if shop, err = repo.ShopRepo.Approve(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error:", err)
		http.Failed(c, code.ErrShopNotExists)
		return
	}
	if shop.Status == shoppingModel.ShopStatusPass {
		go s.addDefaultFriends(req.OperationID, shop.CreatorId)
	}
	http.Success(c)
}

func (s *shoppingService) Update(ctx *gin.Context) {
	var (
		err error
		req shoppingModel.UpdateShopReq
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
	if _, err = shoppingRepo.ShopRepo.UpdateShopByID(req); err != nil {
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

func (s *shoppingService) addDefaultFriends(operationId, userId string) {
	var (
		err error
		fin *configModel.DefaultFriend
	)
	go func() {
		if fin, err = configRepo.DefaultFriendRepo.RandFriend(); err != nil {
			logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" user addDefaultFriends, error: %v", err))
			return
		}
		if err = friendUseCase.FriendUseCase.AddFriend(operationId, userId, fin.UserId, "", fin.GreetMsg, false); err != nil {
			logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" user AddFriend, error: %v", err))
			return
		}
	}()
}
