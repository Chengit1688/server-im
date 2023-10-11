package usecase

import (
	"fmt"
	groupModel "im/internal/api/group/model"
	groupRepo "im/internal/api/group/repo"
	groupUseCase "im/internal/api/group/usecase"
	userModel "im/internal/api/user/model"
	configModel "im/internal/cms_api/config/model"
	"im/pkg/common/constant"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
)

func (s *permissionUseCase) CheckCreateGroupPermission(operationID string, userID string, lang string) (err error) {
	var (
		user   *userModel.UserBaseInfo
		params *configModel.ParameterConfigResp
	)

	if user, params, err = s.GetUserPermissionInfo(operationID, userID, lang); err != nil {
		return
	}

	switch user.IsPrivilege {
	case 1:

	case 2:

		if params.IsMemberAddGroup == constant.SwitchOff {
			err = response.GetError(response.ErrGroupNormalNotCreate, lang)
			return
		}
	}

	var count int64
	if count, err = groupRepo.GroupMemberRepo.GroupCount(userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db group count error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if params.CreateGroupLimit > 0 && count >= params.CreateGroupLimit {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group count max, count: %d, limit count: %d", count, params.CreateGroupLimit))
		err = response.GetError(response.ErrGroupIsMax, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckJoinGroupPermission(operationID string, userID string, groupID string, lang string) (needVerify bool, err error) {
	var (
		user   *userModel.UserBaseInfo
		params *configModel.ParameterConfigResp
	)

	if user, params, err = s.GetUserPermissionInfo(operationID, userID, lang); err != nil {
		return
	}

	var group groupModel.GroupInfo
	if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	switch user.IsPrivilege {
	case 1:

	case 2:

		if params.IsNormalJoinGroup == constant.SwitchOff {
			err = response.GetError(response.ErrUserPermissions, lang)
			return
		}

		if group.JoinNeedApply == 1 {
			needVerify = true
		}
	}

	if params.GroupLimit > 0 && group.MembersTotal >= int(params.GroupLimit) {
		err = response.GetError(response.ErrGroupMemberMax, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckGroupInformationPermission(operationID string, userID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if _, member, err = s.GetGroupAndMemberPermissionInfo(operationID, userID, groupID, lang); err != nil {
		return
	}

	switch member.Role {
	case groupModel.RoleTypeUser:
		err = response.GetError(response.ErrNoPermission, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckGroupManagePermission(operationID string, userID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if _, member, err = s.GetGroupAndMemberPermissionInfo(operationID, userID, groupID, lang); err != nil {
		return
	}

	switch member.Role {
	case groupModel.RoleTypeStaff, groupModel.RoleTypeUser:
		err = response.GetError(response.ErrNoPermission, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckGroupMutePermission(operationID string, userID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	switch member.Role {
	case groupModel.RoleTypeOwner, groupModel.RoleTypeAdmin:

		return

	default:
		err = response.GetError(response.ErrUserPermissions, lang)
	}
	return
}

func (s *permissionUseCase) CheckGroupMemberMutePermission(operationID string, userID string, memberID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	switch member.Role {
	case groupModel.RoleTypeOwner:
		return

	case groupModel.RoleTypeAdmin:
		var groupMember *groupModel.GroupMember
		if groupMember, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, memberID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("admin mute error, get group member error, error: %v", err))
			err = response.GetError(response.ErrGroupNotMember, lang)
			return
		}

		if groupMember.Role == groupModel.RoleTypeStaff || groupMember.Role == groupModel.RoleTypeUser {
			return
		}

	case groupModel.RoleTypeStaff:
		var groupMember *groupModel.GroupMember
		if groupMember, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, memberID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("admin mute error, get group member error, error: %v", err))
			err = response.GetError(response.ErrGroupNotMember, lang)
			return
		}

		if groupMember.Role == groupModel.RoleTypeUser {
			return
		}
	}

	err = response.GetError(response.ErrUserPermissions, lang)
	return
}

func (s *permissionUseCase) CheckGroupMemberInvitePermission(operationID string, userID string, groupID string, count int, lang string) (err error) {
	var (
		params *configModel.ParameterConfigResp
		member *groupModel.GroupMember
	)

	var group groupModel.GroupInfo
	if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if member, params, err = s.GetGroupMemberPermissionInfo(operationID, userID, groupID, lang); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	if member.Role == groupModel.RoleTypeUser {
		err = response.GetError(response.ErrUserPermissions, lang)
		return
	}

	if params.GroupLimit > 0 && group.MembersTotal+count >= int(params.GroupLimit) {
		err = response.GetError(response.ErrGroupMemberMax, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckGroupMemberKickPermission(operationID string, userID string, memberID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	switch member.Role {
	case groupModel.RoleTypeOwner:
		return

	case groupModel.RoleTypeAdmin:
		var groupMember *groupModel.GroupMember
		if groupMember, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, memberID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("admin kick error, get group member error, error: %v", err))
			err = response.GetError(response.ErrGroupNotMember, lang)
			return
		}

		if groupMember.Role == groupModel.RoleTypeStaff || groupMember.Role == groupModel.RoleTypeUser {
			return
		}

	case groupModel.RoleTypeStaff:
		var groupMember *groupModel.GroupMember
		if groupMember, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, memberID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("staff kick error, get group member error, error: %v", err))
			err = response.GetError(response.ErrGroupNotMember, lang)
			return
		}

		if groupMember.Role == groupModel.RoleTypeUser {
			return
		}
	}

	err = response.GetError(response.ErrUserPermissions, lang)
	return
}

func (s *permissionUseCase) CheckGroupRemovePermission(operationID string, userID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember

	if _, member, err = s.GetGroupAndMemberPermissionInfo(operationID, userID, groupID, lang); err != nil {
		return
	}

	if member.Role == groupModel.RoleTypeOwner {
		return
	}

	err = response.GetError(response.ErrNoPermission, lang)
	return
}

func (s *permissionUseCase) CheckGroupQuitPermission(operationID string, userID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	switch member.Role {
	case groupModel.RoleTypeOwner:
		err = response.GetError(response.ErrOwnerCanNotQuit, lang)
		return

	case groupModel.RoleTypeAdmin, groupModel.RoleTypeStaff:
		return

	case groupModel.RoleTypeUser:
		var group groupModel.GroupInfo
		if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group error, error: %v", err))
			err = response.GetError(response.ErrGroupNotExist, lang)
			return
		}

		if group.BanRemoveByNormalMember == 1 {
			err = response.GetError(response.ErrCloseQuit, lang)
			return
		}
	}
	return
}

func (s *permissionUseCase) CheckGroupSetOwnerPermission(operationID string, userID string, memberID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	if member.Role != groupModel.RoleTypeOwner {
		err = response.GetError(response.ErrNoPermission, lang)
		return
	}

	var groupMember *groupModel.GroupMember
	if groupMember, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, memberID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get target group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	if groupMember.Role != groupModel.RoleTypeAdmin {
		err = response.GetError(response.ErrSetOwnerOnlyAdmin, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckGroupSetAdminPermission(operationID string, userID string, memberID string, groupID string, lang string) (err error) {
	var member *groupModel.GroupMember
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	if member.Role != groupModel.RoleTypeOwner {
		err = response.GetError(response.ErrNoPermission, lang)
		return
	}

	if _, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, memberID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get target group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckGroupBaseInfo(operationID string, userID string, groupID string, lang string) (err error) {

	var group groupModel.GroupInfo
	if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrDB, lang)
	}

	if group.Status == 2 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrGroupNotExist, lang)
	}

	if !groupUseCase.GroupMemberUseCase.CheckMember(groupID, userID) {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check member error, error: %v", err))
		return response.GetError(response.ErrGroupNotMember, lang)
	}
	return
}

func (s *permissionUseCase) CheckFace2FaceMember() (err error) {

	return
}
