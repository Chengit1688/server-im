package service

import (
	"gorm.io/gorm"
	"im/config"
	configRepo "im/internal/cms_api/config/repo"
	discoverModel "im/internal/cms_api/discover/model"
	discoverRepo "im/internal/cms_api/discover/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var DiscoverService = new(discoverService)

type discoverService struct{}

func (s *discoverService) GetDiscoverInfo(c *gin.Context) {
	req := new(discoverModel.GetDiscoverInfoReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	discovers, err := discoverRepo.DiscoverRepo.GetDiscovers()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	var ret []discoverModel.GetDiscoverInfoResp
	util.CopyStructFields(&ret, &discovers)
	cfg := config.Config
	for i, item := range discovers {
		ret[i].CreatedAt = item.CreatedAt.Unix()
		if item.Icon == "" {
			ret[i].Icon = cfg.DefaultIcon.DiscoverIcon
		}
	}
	http.Success(c, ret)
}

func (s *discoverService) AddDiscover(c *gin.Context) {
	req := new(discoverModel.AddDiscoverReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	add := new(discoverModel.Discover)
	util.CopyStructFields(&add, &req)
	discover, err := discoverRepo.DiscoverRepo.AddDiscover(*add)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(discoverModel.GetDiscoverInfoResp)
	util.CopyStructFields(&ret, &discover)
	ret.CreatedAt = discover.CreatedAt.Unix()
	http.Success(c, ret)
}

func (s *discoverService) UpdateDiscover(c *gin.Context) {
	req := new(discoverModel.AddDiscoverReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	update := new(discoverModel.Discover)
	util.CopyStructFields(&update, &req)
	discover, err := discoverRepo.DiscoverRepo.UpdateDiscover(util.String2Int(id), *update)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(discoverModel.GetDiscoverInfoResp)
	util.CopyStructFields(&ret, &discover)
	ret.CreatedAt = discover.CreatedAt.Unix()
	http.Success(c, ret)
}

func (s *discoverService) DeleteDiscover(c *gin.Context) {
	req := new(discoverModel.GetDiscoverOpenReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	err = discoverRepo.DiscoverRepo.DeleteDiscover(util.String2Int(id))
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, nil)
}

func (s *discoverService) GetDiscoverOpenStatus(c *gin.Context) {
	req := new(discoverModel.GetDiscoverOpenReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	add := new(discoverModel.Discover)
	util.CopyStructFields(&add, &req)
	status, err := configRepo.ConfigRepo.GetDiscoverIsOpen()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, status)
}

func (s *discoverService) SetDiscoverOpenStatus(c *gin.Context) {
	req := new(discoverModel.SetDiscoverOpenReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configRepo.ConfigRepo.SetDiscoverIsOpen(req.Status)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, nil)
}

func (s *discoverService) AddPrize(c *gin.Context) {
	var (
		req discoverModel.AddPrizeReq
	)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = discoverRepo.DiscoverRepo.AddPrize(req.List)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *discoverService) UpdatePrize(c *gin.Context) {
	var (
		req discoverModel.UpdatePrizeReq
	)
	err := c.ShouldBindJSON(&req)
	if err != nil || req.ID == 0 {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = discoverRepo.DiscoverRepo.UpdatePrize(req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *discoverService) DeletePrize(c *gin.Context) {
	var (
		req discoverModel.DeletePrizeReq
	)
	err := c.ShouldBindJSON(&req)
	if err != nil || req.ID == 0 {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = discoverRepo.DiscoverRepo.DeletePrize(req.ID)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *discoverService) ListPrize(c *gin.Context) {
	req := new(discoverModel.PrizeListReq)
	resp := new(discoverModel.PrizeListResp)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	resp.List, resp.Count, err = discoverRepo.DiscoverRepo.PrizeList(*req)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUnknown)
			return
		}
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *discoverService) RedeemPrizeLog(c *gin.Context) {
	req := new(discoverModel.RedeemPrizeLogReq)
	resp := new(discoverModel.RedeemPrizeLogResp)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	resp.List, resp.Count, err = discoverRepo.DiscoverRepo.RedeemPrizeLog(*req)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUnknown)
			return
		}
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *discoverService) SetRedeemPrize(c *gin.Context) {
	var req discoverModel.SetRedeemPrizeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	prizeLog, err := discoverRepo.DiscoverRepo.FetchRedeemPrizeLog(req.ID)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Status != 0 {
		if prizeLog.IsFictitious == 1 {
			if req.Status != 1 && req.Status != 2 {
				http.Failed(c, code.ErrBadRequest)
				return
			}
		} else {
			if req.Status != 21 && req.Status != 22 {
				http.Failed(c, code.ErrBadRequest)
				return
			}
		}
	}
	_, err = discoverRepo.DiscoverRepo.SetRedeemPrize(req.ID, req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}

	http.Success(c)
}
