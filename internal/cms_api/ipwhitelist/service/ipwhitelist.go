package service

import (
	"im/internal/cms_api/ipwhitelist/model"
	"im/internal/cms_api/ipwhitelist/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

var IPWhiteListService = new(ipWhiteistService)

type ipWhiteistService struct{}

func (s *ipWhiteistService) IPWhiteListPaging(c *gin.Context) {
	req := new(model.GetIPListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	ips, count, err := repo.IPWhiteListRepo.Paging(*req)
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

func (s *ipWhiteistService) IPWhiteListAdd(c *gin.Context) {
	req := new(model.AddIPInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	_, err = repo.IPWhiteListRepo.GetByIP(req.IP)
	if err == nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "IP白名单已存在，不可重复添加")
		http.Failed(c, code.ErrIPWhiteListExist)
		return
	}
	add := new(model.IPWhiteList)
	util.CopyStructFields(&add, &req)
	add.CreatedAt = time.Now().Unix()
	add.UserID = c.GetString("o_user_id")
	ip, err := repo.IPWhiteListRepo.Add(*add)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.AddIPInfoResp)
	util.CopyStructFields(&ret, &ip)
	http.Success(c, ret)
}

func (s *ipWhiteistService) IPWhiteListUpdate(c *gin.Context) {
	req := new(model.UpdateIPInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	_, err = repo.IPWhiteListRepo.Get(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	have, err := repo.IPWhiteListRepo.GetByIP(req.IP)
	if err == nil && int(have.ID) != util.String2Int(id) {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "IP白名单已存在，不可重复添加")
		http.Failed(c, code.ErrIPWhiteListExist)
		return
	}
	update := new(model.IPWhiteList)
	util.CopyStructFields(&update, &req)
	update.UpdatedAt = time.Now().Unix()
	ip, err := repo.IPWhiteListRepo.Update(id, *update)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(model.UpdateIPInfoResp)
	util.CopyStructFields(&ret, &ip)
	http.Success(c, ret)
}

func (s *ipWhiteistService) IPWhiteistDelete(c *gin.Context) {
	req := new(model.DeleteIPInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	_, err = repo.IPWhiteListRepo.Get(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	update := new(model.IPWhiteList)
	util.CopyStructFields(&update, &req)
	err = repo.IPWhiteListRepo.Delete(id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}
