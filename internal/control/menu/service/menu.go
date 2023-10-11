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
	"google.golang.org/grpc/status"
)

var MenuService = new(menuService)

type menuService struct{}

func (s *menuService) HandlerGetMenuConfigTime(c *gin.Context) {
	req := new(model.GetMenuConfigReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	timeString, err := repo.MenuRepo.MenuGetConfigTime()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"))
		http.Failed(c, code.ErrUnknown)
		return
	}

	if req.Timestamp == timeString {
		http.Failed(c, status.Error(201, "配置无变化"))
		return
	} else {
		menus, err := repo.MenuRepo.GetAllMenu()
		if err != nil {
			logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"))
			http.Failed(c, code.ErrUnknown)
			return
		}
		ret := new(model.GetMenuConfigResp)
		ret.Timestamp = timeString
		util.CopyStructFields(&ret.Menus, &menus)
		http.Success(c, ret)
		return
	}
}

func (s *menuService) MenuList(c *gin.Context) {
	req := new(model.MenuListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	roles, count, err := repo.MenuRepo.MenuList(req.Name, req.Title)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.MenuListResp)
	util.CopyStructFields(&ret.List, &roles)
	ret.Count = count
	http.Success(c, ret)
}

func (s *menuService) MenuAdd(c *gin.Context) {
	req := new(model.AddMenuReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id, err := repo.MenuRepo.MenuAdd(*req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db insert error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.AddMenuResp)
	util.CopyStructFields(&ret, &req)
	ret.ID = id
	http.Success(c, ret)
	return
}

func (s *menuService) MenuUpdate(c *gin.Context) {
	req := new(model.UpdateMenuReq)
	id := c.Param("id")
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "params error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = repo.MenuRepo.MenuUpdate(id, *req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(model.UpdateMenuResp)
	util.CopyStructFields(&ret, &req)
	ret.ID = util.StringToInt(id)
	http.Success(c, ret)
}

func (s *menuService) MenuDelete(c *gin.Context) {
	req := new(model.DeleteMenuReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	err = repo.MenuRepo.MenuDelete(id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
}

func (s *menuService) MenuGet(c *gin.Context) {
	req := new(model.DeleteMenuReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	id := c.Param("id")
	menu, err := repo.MenuRepo.MenuGet(id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "db query error"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	menus := []int{}
	for _, item := range *menu.Apis {
		menus = append(menus, int(item.ID))
	}
	ret := new(model.GetMenuResp)
	util.CopyStructFields(&ret, &menu)
	ret.Apis = menus
	http.Success(c, ret)
}
