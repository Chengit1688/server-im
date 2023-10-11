package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	friendUseCase "im/internal/api/friend/usecase"
	shoppingModel "im/internal/api/operator/model"
	shoppingRepo "im/internal/api/operator/repo"
	configModel "im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/internal/cms_api/operator/model"
	"im/internal/cms_api/operator/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var OperatorService = new(operatorService)

type operatorService struct{}

func (s *operatorService) ShopList(c *gin.Context) {
	var (
		err   error
		count int64
		req   model.OperatorListReq
		resp  model.OperatorListResp
		shops []shoppingModel.Operator
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	req.Check()
	if shops, count, err = repo.OperatorRepo.FetchList(req); err != nil {
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

func (s *operatorService) MemberList(ctx *gin.Context) {
	var (
		err  error
		req  model.OperatorMemberListReq
		resp model.OperatorMemberListResp
	)
	lang := ctx.GetHeader("Locale")
	if err = ctx.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	req.Check()
	if _, err = shoppingRepo.OperatorRepo.CheckShop(req.ShopID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(ctx, response.GetError(response.ErrShopNotExists, lang))
		return
	}
	if resp.List, resp.Count, err = repo.OperatorRepo.FetchTeamList(req); err != nil {
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

func (s *operatorService) Approve(c *gin.Context) {
	var (
		err  error
		req  model.OperatorApproveReq
		shop shoppingModel.Operator
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	if shop, err = repo.OperatorRepo.Approve(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error:", err)
		http.Failed(c, code.ErrShopNotExists)
		return
	}
	if shop.Status == shoppingModel.ShopStatusPass {
		go s.addDefaultFriends(req.OperationID, shop.CreatorId)
	}
	http.Success(c)
}

func (s *operatorService) Update(ctx *gin.Context) {
	var (
		err error
		req shoppingModel.UpdateOperatorReq
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
	if _, err = shoppingRepo.OperatorRepo.UpdateShopByID(req); err != nil {
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

func (s *operatorService) addDefaultFriends(operationId, userId string) {
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
