package service

import (
	"fmt"

	apiGroupUseCase "im/internal/api/group/usecase"
	apiUserUseCase "im/internal/api/user/usecase"
	adminModel "im/internal/cms_api/admin/model"
	adminRepo "im/internal/cms_api/admin/repo"
	"im/internal/cms_api/config/model"
	configRepo "im/internal/cms_api/config/repo"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

var InviteCodeService = new(inviteCodeService)

type inviteCodeService struct{}

func (s *inviteCodeService) Add(ctx *gin.Context) {
	var (
		err       error
		req       model.InviteCodeReq
		data      *model.InviteCode
		adminUser *adminModel.Admin
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	operationUserId := ctx.GetString("o_user_id")
	if operationUserId == "" {
		http.Failed(ctx, code.ErrUserPermissions)
		return
	}
	if req.DefaultGroups != "" {
		if !s.checkGroupExists(req.DefaultGroups) {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite checkGroupExists error:", err)
			http.Failed(ctx, code.ErrGroupNotExist)
			return
		}
	}
	if req.DefaultFriends != "" {
		if !s.checkFriendsExists(req.DefaultFriends) {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite checkFriendsExists error:", err)
			http.Failed(ctx, code.ErrFriendNotExist)
			return
		}
	}
	if err = util.Copy(req, &data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite Copy json error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	opt := configRepo.WhereOptionForInvite{
		InviteCode: req.InviteCode,
	}
	if total, _ := configRepo.InviteCode.Exists(opt); total != 0 {
		http.Failed(ctx, code.ErrInviteCodeExist)
		return
	}
	if adminUser, err = adminRepo.AdminRepo.GetByUserID(operationUserId); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite AdminRepo GetByUserID error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	data.OperationUser = adminUser.Username
	if _, err = configRepo.InviteCode.Create(data); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite Create error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}

	http.Success(ctx)
}

func (s *inviteCodeService) checkGroupExists(groups string) bool {
	groupList := strings.Split(groups, ",")
	for _, v := range groupList {
		info, err := apiGroupUseCase.GroupUseCase.GroupInfo(v)
		if err != nil {
			logger.Sugar.Error(util.GetSelfFuncName(), " invite checkGroupExists error:", err)
			return false
		}
		if info.Status != 1 {
			return false
		}
	}
	return true
}

func (s *inviteCodeService) checkFriendsExists(friends string) bool {
	groupList := strings.Split(friends, ",")
	for _, v := range groupList {
		info, err := apiUserUseCase.UserUseCase.GetInfo(v)
		if err != nil {
			logger.Sugar.Error(util.GetSelfFuncName(), " invite checkFriendsExists error:", err)
			return false
		}
		if info.Status != 1 {
			return false
		}
	}
	return true
}

func (s *inviteCodeService) Delete(ctx *gin.Context) {
	var (
		err error
		req model.InviteDeleteReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	opt := configRepo.WhereOptionForInvite{
		Ids: req.Ids,
	}
	if err = configRepo.InviteCode.UpdateById(opt, &model.InviteCode{DeleteStatus: constant.SwitchOn}); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *inviteCodeService) Update(ctx *gin.Context) {
	var (
		err error
		req model.InviteUpdateReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	ep := configRepo.WhereOptionForInvite{
		InviteCode: req.InviteCode,
		Ext:        fmt.Sprintf("id != %d", req.Id),
	}
	if total, _ := configRepo.InviteCode.Exists(ep); total != 0 {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy Exists error:", err)
		http.Failed(ctx, code.ErrInviteCodeExist)
		return
	}
	opt := configRepo.WhereOptionForInvite{
		Id: req.Id,
	}

	if err = configRepo.InviteCode.UpdateInviteById(opt, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update user error: %v", err))
		http.Failed(ctx, code.ErrDB)
		return
	}

	http.Success(ctx)
}

func (s *inviteCodeService) UpdateFriend(ctx *gin.Context) {
	var (
		err error
		req model.InviteUpdateFriendReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	if req.DefaultFriends != "" {
		if !s.checkFriendsExists(req.DefaultFriends) {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite checkFriendsExists error:", err)
			http.Failed(ctx, code.ErrFriendNotExist)
			return
		}
	}

	ep := configRepo.WhereOptionForInvite{
		InviteCode: req.InviteCode,
		Ext:        fmt.Sprintf("id != %d", req.Id),
	}
	if total, _ := configRepo.InviteCode.Exists(ep); total != 0 {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy Exists error:", err)
		http.Failed(ctx, code.ErrInviteCodeExist)
		return
	}
	opt := configRepo.WhereOptionForInvite{
		Id: req.Id,
	}

	updateMap := map[string]interface{}{}
	updateMap["default_friends"] = req.DefaultFriends
	updateMap["friend_index"] = 0

	if err = configRepo.InviteCode.UpdateMapInfo(opt, updateMap); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *inviteCodeService) UpdateGroup(ctx *gin.Context) {
	var (
		err error
		req model.InviteUpdateGroupReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	if req.DefaultGroups != "" {
		if !s.checkGroupExists(req.DefaultGroups) {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite checkGroupExists error:", err)
			http.Failed(ctx, code.ErrGroupNotExist)
			return
		}
	}

	ep := configRepo.WhereOptionForInvite{
		InviteCode: req.InviteCode,
		Ext:        fmt.Sprintf("id != %d", req.Id),
	}
	if total, _ := configRepo.InviteCode.Exists(ep); total != 0 {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy Exists error:", err)
		http.Failed(ctx, code.ErrInviteCodeExist)
		return
	}
	opt := configRepo.WhereOptionForInvite{
		Id: req.Id,
	}

	updateMap := map[string]interface{}{}
	updateMap["default_groups"] = req.DefaultGroups
	updateMap["group_index"] = 0

	if err = configRepo.InviteCode.UpdateMapInfo(opt, updateMap); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}

func (s *inviteCodeService) GetList(ctx *gin.Context) {
	var (
		err error
		req model.InviteListReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}

	list, count, err := configRepo.InviteCode.GetList(req)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite GetList error:", err)
		http.Failed(ctx, code.ErrUnknown)
		return
	}
	ret := new(model.InviteListResp)
	_ = util.CopyStructFields(&ret.List, &list)
	ret.Count = count
	ret.Page = req.Page
	ret.PageSize = req.PageSize

	for i, info := range list {
		groupIDList := strings.Split(info.DefaultGroups, ",")
		for _, groupID := range groupIDList {
			if groupID == "" {
				continue
			}

			groupInfo, _ := apiGroupUseCase.GroupUseCase.GroupInfo(groupID)
			ret.List[i].DefaultGroups = append(ret.List[i].DefaultGroups, groupInfo)
		}

		friendIDList := strings.Split(info.DefaultFriends, ",")
		for _, userID := range friendIDList {
			if userID == "" {
				continue
			}

			userInfo, _ := apiUserUseCase.UserUseCase.GetBaseInfo(userID)
			ret.List[i].DefaultFriends = append(ret.List[i].DefaultFriends, userInfo)
		}
	}

	http.Success(ctx, ret)
}

func (s *inviteCodeService) UpdateStatus(ctx *gin.Context) {
	var (
		err error
		req model.InviteUpdateStatusReq
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite bind json error:", err)
		http.Failed(ctx, code.GetError(err, req))
		return
	}
	opt := configRepo.WhereOptionForInvite{
		Id: req.Id,
	}
	if err = configRepo.InviteCode.UpdateById(opt, &model.InviteCode{Status: req.Status, IsOpenTurn: req.IsOpenTurn}); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " invite UpdateBy error:", err)
		http.Failed(ctx, err)
		return
	}

	http.Success(ctx)
}
