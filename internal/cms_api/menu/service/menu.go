package service

import (
	"im/internal/control/menu/model"
	"im/internal/control/menu/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var MenuService = new(menuService)

type menuService struct{}

func (s *menuService) MenuList(c *gin.Context) {
	req := new(model.MenuListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	roles, count, err := repo.MenuRepo.MenuList(req.Name, req.Title)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.MenuListResp)
	util.CopyStructFields(&ret.List, &roles)
	ret.Count = count
	http.Success(c, ret)
}
