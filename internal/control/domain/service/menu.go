package service

import (
	"fmt"
	"im/internal/control/domain/model"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

var DomainSiteService = new(domainSiteService)

type domainSiteService struct{}

func (s *domainSiteService) AddDomain(c *gin.Context) {
	var (
		req  model.AddDomainReq
		resp model.AddDomainResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	had := model.DomainSite{}
	err = db.Info(&had, model.DomainSite{
		Site:   req.Site,
		Domain: req.Domain,
	})
	if err == nil {
		http.Success(c, resp)
		return
	}
	inData := model.DomainSite{
		Site:   req.Site,
		Domain: req.Domain,
	}
	err = db.Insert(&inData)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	http.Success(c, resp)
}

func (s *domainSiteService) RemoveDomain(c *gin.Context) {
	var (
		req  model.RemoveDomainReq
		resp model.RemoveDomainResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	had := model.DomainSite{}
	err = db.Info(&had, model.DomainSite{
		Site:   req.Site,
		Domain: req.Domain,
	})
	if err != nil {
		http.Success(c, resp)
		return
	}
	inData := model.DomainSite{
		Site:   req.Site,
		Domain: req.Domain,
	}
	err = db.Delete(model.DomainSite{}, &inData)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	http.Success(c, resp)
}

func (s *domainSiteService) DomainList(c *gin.Context) {
	var (
		req  model.DomainListReq
		resp model.DomainListResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	where := map[string]interface{}{}

	if req.Domain != "" {
		where["domain"] = fmt.Sprintf("?%s", req.Domain)
	}
	if req.Site != "" {
		where["site"] = fmt.Sprintf("?%s", req.Site)
	}
	total := int64(0)
	data := []model.DomainSite{}
	db.Find(model.DomainSite{}, where, "id asc", req.Page, req.PageSize, &total, &data)

	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	logger.Sugar.Debug("数据库查出列表", data)
	util.CopyStructFields(&resp.List, data)
	http.Success(c, resp)
}

func (s *domainSiteService) AppDomainList(c *gin.Context) {
	var (
		req  model.AppDomainListReq
		resp model.AppDomainListResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	where := map[string]interface{}{}
	if req.Site != "" {
		where["site"] = req.Site
	}
	total := int64(0)
	data := []model.DomainSite{}
	db.Find(model.DomainSite{}, where, "id desc", 1, 20, &total, &data)

	for _, v := range data {
		resp = append(resp, v.Domain)
	}
	http.Success(c, resp)
}

func (s *domainSiteService) AddWarning(c *gin.Context) {
	var (
		req  model.AddWarningReq
		resp model.AddWarningResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	had := []model.DomainWarning{}
	total := int64(0)
	err = db.Find(model.DomainWarning{}, model.DomainWarning{
		Domain: req.Domain,
	}, "id desc", 1, 1, &total, &had)
	if err == nil && len(had) > 0 && time.Now().Add(-60*time.Minute).Before(had[0].CreatedAt) {
		http.Success(c, resp)
		return
	}
	ip := c.ClientIP()
	db.Insert(&model.DomainWarning{
		Domain: req.Domain,
		Ip:     ip,
	})

	http.Success(c, resp)
}

func (s *domainSiteService) WarningList(c *gin.Context) {
	var (
		req  model.WarningListReq
		resp model.WarningListResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	where := map[string]interface{}{}

	if req.Domain != "" {
		where["domain"] = fmt.Sprintf("?%s", req.Domain)
	}
	total := int64(0)
	data := []model.DomainWarning{}
	db.Find(model.DomainWarning{}, where, "id desc", req.Page, req.PageSize, &total, &data)

	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range data {
		resp.List = append(resp.List, model.WarningInfo{
			ID:        v.ID,
			Domain:    v.Domain,
			Ip:        v.Ip,
			CreatedAt: v.CreatedAt.Unix(),
		})
	}
	http.Success(c, resp)
}
