package service

import (
	"im/internal/cms_api/cmsui/model"
	"im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var CmsUIService = new(cmsuiService)

type cmsuiService struct{}

func (s *cmsuiService) SetCmsSiteUI(c *gin.Context) {
	req := new(model.SetCmsUIDataNormalReq)
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	switch req.Type {
	case 1:
		req2 := new(model.SetCmsUIDataSiteNameReq)
		err := c.ShouldBindBodyWith(&req2, binding.JSON)
		if err != nil {
			logger.Sugar.Errorw(req2.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrBadRequest)
			return
		}
		err = repo.ConfigRepo.SetCmsSiteName(req2.Value)
		if err != nil {
			logger.Sugar.Errorw(req2.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrDB)
			return
		}
		http.Success(c)
		return
	case 2:
		err = repo.ConfigRepo.SetCmsLoginIcon(req.Value)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrDB)
			return
		}
		http.Success(c)
		return
	case 3:
		err = repo.ConfigRepo.SetCmsLoginBackend(req.Value)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrDB)
			return
		}
		http.Success(c)
		return
	case 4:
		err = repo.ConfigRepo.SetCmsPageIcon(req.Value)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrDB)
			return
		}
		http.Success(c)
		return
	case 5:
		err = repo.ConfigRepo.SetCmsMenuIcon(req.Value)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrDB)
			return
		}
		http.Success(c)
		return
	}
	http.Success(c)
}
