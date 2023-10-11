package service

import (
	"im/internal/cms_api/ipblacklist/model"
	"im/internal/cms_api/ipblacklist/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

var IPBlackListService = new(ipblacklistService)

type ipblacklistService struct{}

func (s *ipblacklistService) IPBlackListPaging(c *gin.Context) {
	req := new(model.GetIPListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	ips, count, err := repo.IPBlackListRepo.Paging(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.GetIPListResp)
	util.CopyStructFields(&ret.List, &ips)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *ipblacklistService) IPBlackListAdd(c *gin.Context) {
	req := new(model.AddIPInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	_, err = repo.IPBlackListRepo.GetByIP(req.IP)
	if err == nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "IP黑名单已存在，不可重复添加")
		http.Failed(c, code.ErrIPBlackListExist)
		return
	}
	add := new(model.IPBlackList)
	util.CopyStructFields(&add, &req)
	ip, err := repo.IPBlackListRepo.Add(*add)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.AddIPInfoResp)
	util.CopyStructFields(&ret, &ip)
	http.Success(c, ret)
}

func (s *ipblacklistService) IPBlackListUpdate(c *gin.Context) {
	req := new(model.UpdateIPInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	_, err = repo.IPBlackListRepo.Get(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	have, err := repo.IPBlackListRepo.GetByIP(req.IP)
	if err == nil && int(have.ID) != util.String2Int(id) {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "此IP黑名单已存在，无法修改")
		http.Failed(c, code.ErrIPBlackListExist)
		return
	}
	update := new(model.IPBlackList)
	util.CopyStructFields(&update, &req)
	ip, err := repo.IPBlackListRepo.Update(id, *update)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.UpdateIPInfoResp)
	util.CopyStructFields(&ret, &ip)
	http.Success(c, ret)
}

func (s *ipblacklistService) IPBlackListDelete(c *gin.Context) {
	req := new(model.DeleteIPInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	_, err = repo.IPBlackListRepo.Get(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	update := new(model.IPBlackList)
	util.CopyStructFields(&update, &req)
	err = repo.IPBlackListRepo.Delete(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *ipblacklistService) RemoveInBatch(c *gin.Context) {
	req := new(model.DeleteIPBatchReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if len(req.Ips) == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "IP列表为空")
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = repo.IPBlackListRepo.DeleteInBatch(req.Ips)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *ipblacklistService) AddInBatch(c *gin.Context) {
	req := new(model.AddInBatchReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = repo.IPBlackListRepo.AddInBatch(req.Ips)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
	return
}
