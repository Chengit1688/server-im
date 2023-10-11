package service

import (
	"fmt"
	"im/config"
	groupModel "im/internal/api/group/model"
	groupUseCase "im/internal/api/group/usecase"
	operatorModel "im/internal/api/operator/model"
	shoppingModel "im/internal/api/shopping/model"
	apiUserModel "im/internal/api/user/model"
	apiUserRepo "im/internal/api/user/repo"
	apiUserUseCase "im/internal/api/user/usecase"
	"im/internal/cms_api/config/model"
	configUseCase "im/internal/cms_api/config/usecase"
	operatorRepo "im/internal/cms_api/operator/repo"
	shoppingRepo "im/internal/cms_api/shopping/repo"
	userModel "im/internal/cms_api/user/model"
	userRepo "im/internal/cms_api/user/repo"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/common/constant"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	http2 "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var UserService = new(userService)

type userService struct{}

func (s *userService) UserList(c *gin.Context) {
	req := new(userModel.UserListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	users, count, err := userRepo.UserRepo.UserPaging(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(userModel.UserListResp)
	util.CopyStructFields(&ret.List, &users)
	for index := range ret.List {
		if ret.List[index].LoginIp != "" {
			region, err := util.QueryIpRegion(ret.List[index].LoginIp)
			if err == nil {
				ret.List[index].IPInfo = region
			}
		}
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *userService) RealNameList(c *gin.Context) {
	req := new(userModel.RealNameListReq)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	users, count, err := userRepo.UserRepo.RealNameListPaging(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(userModel.RealNameListResp)
	_ = util.CopyStructFields(&ret.List, &users)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *userService) FindUserToGroupList(c *gin.Context) {
	req := new(userModel.UserForGroupListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	condition := userModel.UserListReq{
		NickName:    req.SearchKey,
		PhoneNumber: req.SearchKey,
		UserID:      req.SearchKey,
		Status:      constant.UserStatusNormal,
		Pagination:  req.Pagination,
	}
	users, count, err := userRepo.UserRepo.UserGroupPaging(condition)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(userModel.UserListResp)
	_ = util.CopyStructFields(&ret.List, &users)
	for index := range ret.List {
		if ret.List[index].LoginIp != "" {
			region, err := util.QueryIpRegion(ret.List[index].LoginIp)
			if err == nil {
				ret.List[index].IPInfo = region
			}
		}
	}
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}

func (s *userService) UserBatchAdd(c *gin.Context) {
	req := new(userModel.UserAddBatchReq)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	checkSameAccountMap := map[string]bool{}

	checkSameNicknameMap := map[string]bool{}
	for _, checkAccount := range req.Users {
		if _, ok := checkSameAccountMap[checkAccount.Account]; ok {
			http.Failed(c, code.ErrSameAccount)
			return
		}
		checkSameAccountMap[checkAccount.Account] = true

		if _, ok := checkSameNicknameMap[checkAccount.NickName]; ok {
			http.Failed(c, code.ErrSameNickname)
			return
		}
		checkSameNicknameMap[checkAccount.NickName] = true
	}

	opts := new(apiUserRepo.WhereOption)
	for index, check := range req.Users {
		opts.Account = strings.ToLower(check.Account)
		exist := apiUserRepo.UserRepo.OrExists(*opts)
		if exist {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "账号已存在", "account", opts.Account)
			http.Failed(c, code.ErrAccountExist)
			return
		}
		req.Users[index].Account = strings.ToLower(check.Account)
	}

	optsNickname := new(apiUserRepo.WhereOption)
	for _, check := range req.Users {
		optsNickname.NickName = check.NickName
		_, err = apiUserRepo.UserRepo.GetInfo(*optsNickname)
		if err == nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "昵称已存在", "nickname", optsNickname.NickName)
			http.Failed(c, code.ErrNickNameUsed)
			return
		}
	}
	users, err := userRepo.UserRepo.UserBatchAdd(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(userModel.UserAddBatchResp)
	util.CopyStructFields(&ret.Users, &users)
	for index := range ret.Users {
		ret.Users[index].Password = "******"
	}
	http.Success(c, ret)
}

func (s *userService) UserDetails(c *gin.Context) {
	req := new(userModel.GetUserDetailsByUserIDReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	opts := new(apiUserRepo.WhereOption)
	opts.UserId = req.UserID
	user, err := apiUserRepo.UserRepo.GetByUserID(*opts)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	ret := new(userModel.GetUserDetailsByUserIDResp)
	util.CopyStructFields(&ret, &user)
	ret.RegisterTime = user.CreatedAt

	ret.Online = 2
	if apiUserUseCase.UserUseCase.IsOnline(req.UserID) {
		ret.Online = 1
	}
	http.Success(c, ret)
}

func (s *userService) UserInfoUpdate(c *gin.Context) {
	req := new(userModel.UpdateUserInfoReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	var userInfo *apiUserModel.User
	opts := new(apiUserRepo.WhereOption)
	opts.UserId = req.UserID
	userInfo, err = apiUserRepo.UserRepo.GetByUserID(*opts)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	optsNickname := new(apiUserRepo.WhereOption)

	optsNickname.NickName = *req.NickName
	dbUser, err := apiUserRepo.UserRepo.GetInfo(*optsNickname)
	if err == nil {
		if dbUser.UserID != *&req.UserID {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "昵称已存在", "nickname", optsNickname.NickName)
			http.Failed(c, code.ErrNickNameUsed)
			return
		}
	}
	if req.Age == nil {
		*req.Age = int64(18)
	}
	err = userRepo.UserRepo.UserInfoUpdate(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	user, err := apiUserRepo.UserRepo.GetByUserID(*opts)
	ret := new(userModel.GetUserDetailsByUserIDResp)
	util.CopyStructFields(&ret, &user)
	ret.RegisterTime = user.CreatedAt

	ret.Online = 2
	if apiUserUseCase.UserUseCase.IsOnline(req.UserID) {
		ret.Online = 1
	}

	apiUserRepo.UserCache.DelUserInfoOnCache(userInfo.UserID)
	updateU := apiUserModel.UserBaseInfo{
		UserId:      user.UserID,
		Account:     user.Account,
		FaceURL:     user.FaceURL,
		BigFaceURL:  user.BigFaceURL,
		Gender:      user.Gender,
		NickName:    user.NickName,
		Signatures:  user.Signatures,
		Age:         user.Age,
		IsPrivilege: user.IsPrivilege,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserInfoPush, updateU, userInfo.UserID); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " SendMessageToOnlineUsers error:", err)
	}
	s.runRoutineByUserId(req.OperationID, userInfo.UserID)
	http.Success(c, ret)
}

func (s *userService) runRoutineByUserId(OperationID, userId string) {
	rm := []apiUserUseCase.FuncForUserId{
		groupUseCase.GroupMemberUseCase.UserInfoChange,
	}
	apiUserUseCase.UserUseCase.DoRoutineByUserId(OperationID, userId, rm...)
}

func (s *userService) FreezeUser(c *gin.Context) {
	req := new(userModel.FreezeUserReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if len(req.UserID) == 0 && len(req.UserIDs) == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "需要传递用户ID")
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if len(req.UserID) != 0 {
		opts := new(apiUserRepo.WhereOption)
		opts.UserId = req.UserID
		_, err = apiUserRepo.UserRepo.GetByUserID(*opts)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUserIdNotExist)
			return
		}
		err = userRepo.UserRepo.FreezeUser(req.UserID)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrDB)
			return
		}

		mqtt.SendMessageToUsers(req.OperationID, common.FreezeUserPush, &userModel.FreezePushMessage{Msg: "用户冻结"}, []string{req.UserID}...)
	}
	if len(req.UserIDs) != 0 {
		for _, UserID := range req.UserIDs {
			opts := new(apiUserRepo.WhereOption)
			opts.UserId = UserID
			_, err = apiUserRepo.UserRepo.GetByUserID(*opts)
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
				http.Failed(c, code.ErrUserIdNotExist)
				return
			}
			err = userRepo.UserRepo.FreezeUser(UserID)
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
				http.Failed(c, code.ErrDB)
				return
			}
		}

		mqtt.SendMessageToUsers(req.OperationID, common.FreezeUserPush, &userModel.FreezePushMessage{Msg: "用户冻结"}, req.UserIDs...)
	}

	http.Success(c, nil)
}

func (s *userService) UnFreezeUser(c *gin.Context) {
	req := new(userModel.UnFreezeUserReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if len(req.UserID) == 0 && len(req.UserIDs) == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "需要传递用户ID")
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if len(req.UserID) != 0 {
		opts := new(apiUserRepo.WhereOption)
		opts.UserId = req.UserID
		_, err = apiUserRepo.UserRepo.GetByUserID(*opts)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUserIdNotExist)
			return
		}
		err = userRepo.UserRepo.UnFreezeUser(req.UserID)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrDB)
			return
		}
	}
	if len(req.UserIDs) != 0 {
		for _, UserID := range req.UserIDs {
			opts := new(apiUserRepo.WhereOption)
			opts.UserId = UserID
			_, err = apiUserRepo.UserRepo.GetByUserID(*opts)
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
				http.Failed(c, code.ErrUserIdNotExist)
				return
			}
			err = userRepo.UserRepo.UnFreezeUser(UserID)
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
				http.Failed(c, code.ErrDB)
				return
			}
		}
	}
	http.Success(c, nil)
}

func (s *userService) SetUserPassword(c *gin.Context) {
	req := new(userModel.SetUserPasswordReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	opts := new(apiUserRepo.WhereOption)
	opts.UserId = req.UserID
	_, err = apiUserRepo.UserRepo.GetByUserID(*opts)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	err = userRepo.UserRepo.SetUserPassword(req.UserID, req.Password)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c, nil)
}

func (s *userService) SignLogList(c *gin.Context) {
	var (
		err   error
		count int64
		req   model.SignLogListReq
		resp  model.SignLogListResp
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " sign log bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	if resp.List, count, err = userRepo.SignLogRepo.GetList(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " sign log GetList error:", err)
		http.Failed(c, code.ErrBadSignLog)
		return
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(c, resp)
}

func (s *userService) AgentLevel(c *gin.Context) {
	var (
		err   error
		count int64
		req   shoppingModel.AgentLevelListReq
		resp  shoppingModel.AgentLevelListResp
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	req.Check()
	if resp.List, count, err = shoppingRepo.ShopRepo.FetchTeamListByInviteUserId(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " error", err)
		http.Failed(c, code.ErrBadSignLog)
		return
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(c, resp)
}

func (s *userService) OperatorLevel(c *gin.Context) {
	var (
		err   error
		count int64
		req   operatorModel.OperatorAgentLevelListReq
		resp  operatorModel.OperatorAgentLevelListResp
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	req.Check()
	if resp.List, count, err = operatorRepo.OperatorRepo.FetchTeamListByInviteUserId(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " error", err)
		http.Failed(c, code.ErrBadSignLog)
		return
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(c, resp)
}

func (s *userService) RealNameAuth(c *gin.Context) {
	var (
		err error
		req userModel.RealNameAuthReq
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	if err = userRepo.UserRepo.UpdateRealName(req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	http.Success(c)
}

func (s *userService) CustomerUserList(c *gin.Context) {
	var (
		err   error
		count int64
		req   model.PrivilegeUserReq
		resp  model.PrivilegeUserListResp
		users []apiUserModel.User
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	userListReq := userModel.UserListReq{
		IsCustomer: constant.SwitchOn,

		Account:    req.Account,
		NickName:   req.NickName,
		UserID:     req.UserId,
		Pagination: req.Pagination,
	}
	if users, count, err = userRepo.UserRepo.UserPaging(userListReq); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege GetList error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}
	if err = util.CopyStructFields(&resp.List, users); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " copy json error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(c, resp)
}

func (s *userService) CustomerUserAdd(c *gin.Context) {
	var (
		err  error
		req  model.PrivilegeUserReq
		user *apiUserModel.User
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	opt := apiUserRepo.WhereOption{
		Account:  req.Account,
		NickName: req.NickName,
		UserId:   req.UserId,
	}
	if user, err = apiUserRepo.UserRepo.GetInfo(opt); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " customer add error:", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	if user.IsCustomer == constant.SwitchOn {
		http.Failed(c, code.ErrCustomerUserExist)
		return
	}
	opt = apiUserRepo.WhereOption{
		Id: user.ID,
	}
	userData := &apiUserModel.User{
		IsCustomer: constant.SwitchOn,
	}

	if _, err = apiUserRepo.UserRepo.UpdateBy(opt, userData); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " customer add error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}

	http.Success(c)
}

func (s *userService) CustomerUserRemove(c *gin.Context) {
	var (
		err  error
		req  model.PrivilegeUserReq
		user *apiUserModel.User
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	opt := apiUserRepo.WhereOption{
		Account:  req.Account,
		NickName: req.NickName,
		UserId:   req.UserId,
	}
	if user, err = apiUserRepo.UserRepo.GetInfo(opt); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " customer add error:", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	opt = apiUserRepo.WhereOption{
		Id: user.ID,
	}
	userData := &apiUserModel.User{
		IsCustomer: constant.SwitchOff,
	}
	if _, err = apiUserRepo.UserRepo.UpdateBy(opt, userData); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege remove error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}

	http.Success(c)
}

func (s *userService) PrivilegeUserList(c *gin.Context) {
	var (
		err   error
		count int64
		req   model.PrivilegeUserReq
		resp  model.PrivilegeUserListResp
		users []apiUserModel.User
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	userListReq := userModel.UserListReq{
		IsPrivilege: constant.SwitchOn,

		Account:    req.Account,
		NickName:   req.NickName,
		UserID:     req.UserId,
		Pagination: req.Pagination,
	}
	if users, count, err = userRepo.UserRepo.UserPaging(userListReq); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege GetList error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}
	if err = util.CopyStructFields(&resp.List, users); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " copy json error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}
	freezeStatus, err := configUseCase.ConfigUseCase.GetPrivilegeUserFreezeIsOpen()
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " GetPrivilegeUserFreezeIsOpen error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	resp.IsFreeze = freezeStatus
	resp.Count = count
	resp.PageSize = req.PageSize
	resp.Page = req.Page

	http.Success(c, resp)
}

func (s *userService) PrivilegeUserAdd(c *gin.Context) {
	var (
		err  error
		req  model.PrivilegeUserReq
		user *apiUserModel.User
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	opt := apiUserRepo.WhereOption{
		Account:  req.Account,
		NickName: req.NickName,
		UserId:   req.UserId,
	}
	if user, err = apiUserRepo.UserRepo.GetInfo(opt); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege add error:", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	if user.IsPrivilege == constant.SwitchOn {
		http.Failed(c, code.ErrPrivilegeUserExist)
		return
	}
	opt = apiUserRepo.WhereOption{
		Id: user.ID,
	}
	userData := &apiUserModel.User{
		IsPrivilege: constant.SwitchOn,
	}

	freezeStatus, err := configUseCase.ConfigUseCase.GetPrivilegeUserFreezeIsOpen()
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " GetPrivilegeUserFreezeIsOpen error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	if freezeStatus == constant.SwitchOn {
		userData.Status = 2

		mqtt.SendMessageToUsers(req.OperationID, common.FreezeUserPush, &userModel.FreezePushMessage{Msg: "用户冻结"}, []string{req.UserId}...)
	}
	if _, err = apiUserRepo.UserRepo.UpdateBy(opt, userData); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege add error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}

	groupUseCase.GroupMemberUseCase.UserRoleChange(user.UserID, groupModel.RoleTypeUser, groupModel.RoleTypeStaff)

	apiUserRepo.UserCache.DelUserInfoOnCache(user.UserID)
	updateU := apiUserModel.UserBaseInfo{
		UserId:      user.UserID,
		Account:     user.Account,
		FaceURL:     user.FaceURL,
		BigFaceURL:  user.BigFaceURL,
		Gender:      user.Gender,
		NickName:    user.NickName,
		Signatures:  user.Signatures,
		Age:         user.Age,
		IsPrivilege: constant.SwitchOn,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserInfoPush, updateU, user.UserID); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " SendMessageToOnlineUsers error:", err)
	}

	http.Success(c)
}

func (s *userService) PrivilegeUserRemove(c *gin.Context) {
	var (
		err  error
		req  model.PrivilegeUserReq
		user *apiUserModel.User
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.GetError(err, req))
		return
	}
	opt := apiUserRepo.WhereOption{
		Account:  req.Account,
		NickName: req.NickName,
		UserId:   req.UserId,
	}
	if user, err = apiUserRepo.UserRepo.GetInfo(opt); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege add error:", err)
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	opt = apiUserRepo.WhereOption{
		Id: user.ID,
	}
	userData := &apiUserModel.User{
		IsPrivilege: constant.SwitchOff,
	}
	if _, err = apiUserRepo.UserRepo.UpdateBy(opt, userData); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " privilege remove error:", err)
		http.Failed(c, code.ErrFailRequest)
		return
	}

	groupUseCase.GroupMemberUseCase.UserRoleChange(user.UserID, groupModel.RoleTypeStaff, groupModel.RoleTypeUser)

	apiUserRepo.UserCache.DelUserInfoOnCache(user.UserID)
	updateU := apiUserModel.UserBaseInfo{
		UserId:      user.UserID,
		Account:     user.Account,
		FaceURL:     user.FaceURL,
		BigFaceURL:  user.BigFaceURL,
		Gender:      user.Gender,
		NickName:    user.NickName,
		Signatures:  user.Signatures,
		Age:         user.Age,
		IsPrivilege: constant.SwitchOff,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
	}

	if err = mqtt.SendMessageToUsers(req.OperationID, common.UserInfoPush, updateU, user.UserID); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " SendMessageToOnlineUsers error:", err)
	}

	http.Success(c)
}

func (s *userService) DisabledManagermentUser(c *gin.Context) {
	var (
		err   error
		req   userModel.DMUserListReq
		resp  userModel.DMUserListResp
		users []userModel.DMUserListItemResp
		count int64
	)
	err = c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	users, count, err = userRepo.UserRepo.DisabledManagermentUser(req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	resp.List = users
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *userService) DisabledManagermentDevice(c *gin.Context) {
	var (
		err     error
		req     userModel.DMDeviceListReq
		resp    userModel.DMDeviceListResp
		devices []userModel.DMDeviceListItemResp
		count   int64
	)
	err = c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	devices, count, err = userRepo.UserRepo.DisabledManagermentDevice(req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	resp.List = devices
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *userService) DisabledManagermentIP(c *gin.Context) {
	var (
		err   error
		req   userModel.DMIPListReq
		resp  userModel.DMIPListResp
		ips   []userModel.DMIPListItemResp
		count int64
	)
	err = c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	ips, count, err = userRepo.UserRepo.DisabledManagermentIP(req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	for index, ip := range ips {
		ipInfo, _ := util.QueryIpRegion(ip.IP)
		ips[index].IPInfo = ipInfo
	}
	resp.List = ips
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *userService) DMDeviceLock(c *gin.Context) {
	var (
		err error
		req userModel.DeviceLockReq
	)
	err = c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = userRepo.UserRepo.DMDeviceLock(req.DeviceIDs)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *userService) DMDeviceUnLock(c *gin.Context) {
	var (
		err error
		req userModel.DeviceLockReq
	)
	err = c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = userRepo.UserRepo.DMDeviceUnLock(req.DeviceIDs)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	http.Success(c)
}

func (s *userService) UserListExport(c *gin.Context) {
	req := new(userModel.UserListReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	users, err := userRepo.UserRepo.UserExport(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(userModel.UserListResp)
	util.CopyStructFields(&ret.List, &users)
	for index := range ret.List {
		if ret.List[index].LoginIp != "" {
			region, err := util.QueryIpRegion(ret.List[index].LoginIp)
			if err == nil {
				ret.List[index].IPInfo = region
			}
		}
	}
	cfg := config.Config.Minio

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		}
	}()

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}

	err = f.SetColWidth("Sheet1", "A", "K", 20)

	sheetHeader := []interface{}{"头像", "ID", "账号", "昵称", "手机号", "注册时间", "最后登录IP", "IP归属地", "余额", "状态", "邀请码"}
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	f.SetSheetRow("Sheet1", cell, &sheetHeader)
	var row []interface{}
	var faceUrl string
	var RegisterTime string
	var Status string
	for index := range ret.List {
		cell, err := excelize.CoordinatesToCellName(1, index+2)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, code.ErrUnknown)
			return
		}
		if strings.HasPrefix(ret.List[index].FaceURL, "http") || ret.List[index].FaceURL == "" {
			faceUrl = ret.List[index].FaceURL
		} else {
			faceUrl = cfg.OssPrefix + ret.List[index].FaceURL
		}

		timeLayout := "2006-01-02 15:04:05"
		RegisterTime = time.Unix(ret.List[index].CreatedAt, 0).Format(timeLayout)

		if ret.List[index].Status == 1 {
			Status = "启用"
		} else {
			Status = "冻结"
		}
		row = []interface{}{faceUrl, ret.List[index].UserID, ret.List[index].Account, ret.List[index].NickName, ret.List[index].PhoneNumber, RegisterTime, ret.List[index].LoginIp, ret.List[index].IPInfo, ret.List[index].Balance, Status, ret.List[index].InviteCode}
		f.SetSheetRow("Sheet1", cell, &row)
	}

	f.SetActiveSheet(index)

	buf, err := f.WriteToBuffer()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	c.Writer.WriteHeader(http2.StatusOK)
	filename := url.QueryEscape("用户列表.xlsx")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write(buf.Bytes())
}

func (s *userService) UserSearch(c *gin.Context) {
	req := new(userModel.UserSearchReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	users, count, err := userRepo.UserRepo.UserSearch(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrUnknown)
		return
	}
	ret := new(userModel.UserSearchResp)
	util.CopyStructFields(&ret.List, &users)
	ret.Count = count
	http.Success(c, ret)
}

func (s *userService) SetPrivilegeUserFreeze(c *gin.Context) {
	req := new(userModel.SetPrivilegeUserFreezeReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = configUseCase.ConfigUseCase.SetPrivilegeUserFreezeIsOpen(req.IsFreeze)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " GetPrivilegeUserFreezeIsOpen error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	var status int
	switch req.IsFreeze {
	case constant.SwitchOn:
		status = 2
	case constant.SwitchOff:
		status = 1
	}
	users, err := userRepo.UserRepo.SetPrivilegeUserStatus(status)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	if req.IsFreeze == constant.SwitchOn {
		var userIds []string
		for _, user := range users {
			userIds = append(userIds, user.UserID)
		}
		mqtt.SendMessageToUsers(req.OperationID, common.FreezeUserPush, &userModel.FreezePushMessage{Msg: "用户冻结"}, userIds...)
	}
	http.Success(c)
}

func (s *userService) LoginHistory(c *gin.Context) {
	req := new(userModel.LoginHistoryReq)
	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	records, count, err := userRepo.UserRepo.LoginHistoryPaging(*req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrDB)
		return
	}
	ret := new(userModel.LoginHistoryResp)
	util.CopyStructFields(&ret.List, &records)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize
	http.Success(c, ret)
}
