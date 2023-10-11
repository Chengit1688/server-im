package service

import (
	roleModel "im/internal/cms_api/role/model"
	roleRepo "im/internal/cms_api/role/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var RoleService = new(roleService)

type roleService struct{}

func (s *roleService) RoleList(c *gin.Context) {
	req := new(roleModel.RoleListReq)
	err := c.ShouldBindQuery(&req)
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	roles, count, err := roleRepo.RoleRepo.RolePaging(*req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(roleModel.RoleListResp)
	util.CopyStructFields(&ret.List, &roles)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
	return
}

func (s *roleService) RoleAdd(c *gin.Context) {
	req := new(roleModel.RoleAddReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	req.Menus = removeDuplicateElement(req.Menus)
	_, err = roleRepo.RoleRepo.RoleGetByKey(req.RoleKey, "")
	if err == nil {
		logger.Sugar.Errorw(req.OperationID, zap.String("func", util.GetSelfFuncName()), zap.String("error", "重复角色 不能添加role_key"))
		http.Failed(c, code.ErrRoleKeyExist)
		return
	}
	_, err = roleRepo.RoleRepo.RoleGetByName(req.RoleName, "")
	if err == nil {
		logger.Sugar.Errorw(req.OperationID, zap.String("func", util.GetSelfFuncName()), zap.String("error", "重复角色 不能添加role_name"))
		http.Failed(c, code.ErrRoleNameExist)
		return
	}

	id, err := roleRepo.RoleRepo.RoleAdd(*req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	res := new(roleModel.RoleAddResp)
	util.CopyStructFields(&res, &req)
	res.ID = id
	http.Success(c, res)
	return
}

func (s *roleService) RoleUpdate(c *gin.Context) {
	req := new(roleModel.UpdateRoleReq)
	id := c.Param("id")
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}

	haveKey, err := roleRepo.RoleRepo.RoleGetByKey(req.RoleKey, id)
	if err == nil && haveKey.ID != util.String2Int(id) {
		logger.Sugar.Errorw(req.OperationID, zap.String("func", util.GetSelfFuncName()), zap.String("error", "重复角色 不能修改role_key"))
		http.Failed(c, code.ErrRoleKeyExist)
		return
	}
	haveName, err := roleRepo.RoleRepo.RoleGetByName(req.RoleName, id)
	if err == nil && haveName.ID != util.String2Int(id) {
		logger.Sugar.Errorw(req.OperationID, zap.String("func", util.GetSelfFuncName()), zap.String("error", "重复角色 不能修改role_name"))
		http.Failed(c, code.ErrRoleNameExist)
		return
	}
	req.Menus = removeDuplicateElement(req.Menus)
	err = roleRepo.RoleRepo.RoleUpdate(id, *req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	res := new(roleModel.UpdateRoleResp)
	util.CopyStructFields(&res, &req)
	res.ID = util.StringToInt(id)
	http.Success(c, res)
}

func (s *roleService) RoleDelete(c *gin.Context) {
	req := new(roleModel.DeleteRoleReq)
	id := c.Param("id")
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = roleRepo.RoleRepo.RoleDelete(id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
}

func (s *roleService) RoleGet(c *gin.Context) {
	req := new(roleModel.GetRoleByIDReq)
	id := c.Param("id")
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	role, err := roleRepo.RoleRepo.RoleGet(id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(roleModel.GetRoleByIDResp)
	util.CopyStructFields(&ret, &role)
	if ret.CmsMenu == nil {
		ret.CmsMenu = []roleModel.GetRoleByIDMenusItemResp{}
	}
	http.Success(c, ret)
}

func removeDuplicateElement(data []int) []int {
	result := make([]int, 0, len(data))
	temp := map[int]struct{}{}
	for _, item := range data {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
