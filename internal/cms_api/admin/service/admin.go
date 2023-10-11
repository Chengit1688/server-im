package service

import (
	"fmt"
	adminModel "im/internal/cms_api/admin/model"
	adminRepo "im/internal/cms_api/admin/repo"
	configRepo "im/internal/cms_api/config/repo"
	roleModel "im/internal/cms_api/role/model"
	roleRepo "im/internal/cms_api/role/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var AdminService = new(adminService)

type adminService struct{}

func (s *adminService) Login(c *gin.Context) {
	req := new(adminModel.LoginReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	user, err := adminRepo.AdminRepo.GetByUsername(req.Username)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUserNotFound)
		return
	}
	status := util.CheckPassword(user.Password, req.Password, user.Salt)
	if !status {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "密码错误"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrWrongPassword)
		return
	}
	role, err := roleRepo.RoleRepo.GetRoleKeyByRoleID(user.Role)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "该用户没有有效角色"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}

	googleCodeIsOpen, err := configRepo.ConfigRepo.GetGoogleCodeIsOpen()
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "获取谷歌验证码开关失败"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	if googleCodeIsOpen == 1 {

		if len(user.Google2fSecretKey) > 0 {

			ok := util.VerifyGoogleCode(user.Google2fSecretKey, req.GoogleCode)
			if !ok {
				logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "谷歌验证码校验失败"), zap.String("operation_id", req.OperationID))
				http.Failed(c, code.ErrBadCode)
				return
			}
		}
	}
	diration := time.Duration(24) * time.Hour * 7
	token, err := util.CmsCreateToken(user.UserID, role.RoleKey, diration)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "创建token失败"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ip := c.ClientIP()
	err = adminRepo.AdminRepo.UpdateLoginInfo(user.UserID, ip)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "更新登录信息失败"), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(adminModel.LoginResp)
	ret.Token = token
	ret.Expire = time.Now().Add(diration).Unix()
	http.Success(c, ret)
	return
}

func (s *adminService) RefreshToken(c *gin.Context) {
	req := new(adminModel.RefreshTokenReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	token := c.Request.Header.Get("token")
	userID, roleKey, exp, _ := util.CmsParseToken(token)
	currentTime := time.Now()
	m, _ := time.ParseDuration("-1h")
	exp_ := currentTime.Add(m)
	ret := new(adminModel.LoginResp)
	ret.Token = token
	if exp < exp_.Unix() {
		ret.Token = token
		ret.Expire = exp
	} else {
		diration := time.Duration(24) * time.Hour * 7
		token, err := util.CmsCreateToken(userID, roleKey, diration)
		if err != nil {
			logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", "创建token失败"), zap.String("operation_id", req.OperationID))
			http.Failed(c, code.ErrUnknown)
			return
		}
		ret.Token = token
		ret.Expire = time.Now().Add(diration).Unix()
	}
	http.Success(c, ret)
	return
}

func (s *adminService) GetInfo(c *gin.Context) {
	userID := c.GetString("o_user_id")
	roleKey := c.GetString("o_role_key")
	data := new(adminModel.GetinfoResp)
	user, err := adminRepo.AdminRepo.GetByUserID(userID)
	if err != nil {
		logger.Sugar.Error("GetInfo", zap.String("func", util.GetSelfFuncName()))
		http.Failed(c, code.ErrUnknown)
		return
	}
	data.NickName = user.Nickname
	data.UserID = userID
	data.Username = user.Username
	data.PhoneNumber = user.PhoneNumber
	if roleKey == "admin" {
		role, err := roleRepo.RoleRepo.GetAdminMenu()
		if err != nil {
			logger.Sugar.Error("GetInfo", zap.String("func", util.GetSelfFuncName()))
			http.Failed(c, code.ErrUnknown)
			return
		}
		menus := makeMenuTreeData(role, 0)
		data.Menus = menus
		http.Success(c, data)
		return
	} else {
		role, err := roleRepo.RoleRepo.GetRoleMenu(roleKey)
		if err != nil {
			logger.Sugar.Error("GetInfo", zap.String("func", util.GetSelfFuncName()))
			http.Failed(c, code.ErrUnknown)
			return
		}
		menus := makeMenuTreeData(*role.CmsMenu, 0)
		data.Menus = menus
		http.Success(c, data)
		return
	}

}

func makeMenuTreeData(cmsMenus []roleModel.CmsMenu, parentId int) (tree []adminModel.GetinfoMenuResp) {

	for _, menu := range cmsMenus {
		var add adminModel.GetinfoMenuResp
		if menu.ParentId == parentId {
			util.CopyStructFields(&add, &menu)
			add.MenuID = int(menu.ID)
			add.MenuName = menu.Name
			add.MenuType = menu.Type
			add.Children = nil

			subMenus := makeMenuTreeData(cmsMenus, add.MenuID)

			for i, subMenu := range subMenus {

				subMenuChildrens := makeMenuTreeData(cmsMenus, subMenu.MenuID)
				subMenu.Children = subMenuChildrens
				subMenus[i] = subMenu
			}
			add.Children = subMenus
			tree = append(tree, add)
		}
	}

	return
}

func (s *adminService) List(c *gin.Context) {
	req := new(adminModel.ListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	users, count, err := adminRepo.AdminRepo.Paging(*req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUserNotFound)
		return
	}
	ret := new(adminModel.ListResp)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	util.CopyStructFields(&ret.List, &users)
	for i, item := range users {
		if !users[i].LastloginTime.IsZero() {
			ret.List[i].LastloginTime = item.LastloginTime.Unix()
		}
	}
	http.Success(c, ret)
	return
}

func (s *adminService) Add(c *gin.Context) {
	req := new(adminModel.AddReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	have, err := adminRepo.AdminRepo.GetByUsername(req.Username)
	if err == nil {
		if have.ID > 0 {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "username already exist. can't create")
			http.Failed(c, code.ErrUserIdExist)
			return
		}
	}
	operation_user_id := c.GetString("o_user_id")
	user, err := adminRepo.AdminRepo.Add(*req, operation_user_id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(adminModel.AddResp)
	util.CopyStructFields(&ret, &user)
	ret.RoleID = int(user.Role)
	http.Success(c, ret)
	return
}

func (s *adminService) UpdateInfo(c *gin.Context) {
	req := new(adminModel.UpdateInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	operation_user_id := c.GetString("o_user_id")
	have, err := adminRepo.AdminRepo.GetByUsername(req.Username)
	if err == nil {
		if have.ID != uint(req.ID) {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "username already exist. can't update")
			http.Failed(c, code.ErrUserIdExist)
			return
		}
	}
	_, err = adminRepo.AdminRepo.UpdateInfo(*req, operation_user_id)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
	return
}

func (s *adminService) UpdatePassword(c *gin.Context) {
	req := new(adminModel.UpdatePasswordReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	operation_user_id := c.GetString("o_user_id")
	_, err = adminRepo.AdminRepo.UpdatePassword(*req, operation_user_id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
	return
}

func (s *adminService) Delete(c *gin.Context) {
	req := new(adminModel.DeleteReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	operation_user_id := c.GetString("o_user_id")
	_, err = adminRepo.AdminRepo.Delete(*req, operation_user_id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrUnknown)
		return
	}
	http.Success(c, nil)
	return
}

func (s *adminService) GetGoogleCodeSecret(c *gin.Context) {
	req := new(adminModel.GetGoogleCodeReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	operation_user_id := c.GetString("o_user_id")
	user, err := adminRepo.AdminRepo.GetGoogleCodeSecret(operation_user_id)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("error", err.Error()), zap.String("operation_id", req.OperationID))
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(adminModel.GetGoogleCodeResp)
	ret.Username = user.Username
	qrcontent := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=IM", user.Username, user.Google2fSecretKey)
	ret.Secret = user.Google2fSecretKey
	ret.Image = util.GetGoogleQRCode(qrcontent)
	http.Success(c, ret)
	return
}
