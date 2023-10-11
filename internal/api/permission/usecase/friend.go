package usecase

import (
	"fmt"
	friendRepo "im/internal/api/friend/repo"
	userModel "im/internal/api/user/model"
	configModel "im/internal/cms_api/config/model"
	"im/pkg/common/constant"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
)

func (s *permissionUseCase) CheckAddFriendPermission(operationID string, userID string, friendID string, lang string) (needVerify bool, err error) {
	var (
		user       *userModel.UserBaseInfo
		friendUser *userModel.UserBaseInfo
		params     *configModel.ParameterConfigResp
	)

	if user, friendUser, params, err = s.GetFriendPermissionInfo(operationID, userID, friendID, lang); err != nil {
		return
	}

	switch user.IsPrivilege {
	case 1:

		if params.IsPrivilegeAddVerify == constant.SwitchOn {
			needVerify = true
		}

	case 2:

		switch friendUser.IsPrivilege {
		case 1:

			if params.IsNormalAddPrivilege == constant.SwitchOff {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("is_normal_add_privilege switch off"))
				err = response.GetError(response.ErrFriendNormalNotAdd, lang)
				return
			}

			if params.IsAddPrivilegeVerify == constant.SwitchOn {
				needVerify = true
			}

		case 2:

			if params.IsMemberAddFriend == constant.SwitchOff {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("is_member_add_friend switch off"))
				err = response.GetError(response.ErrFriendNormalNotAdd, lang)
				return
			}

			if params.IsAddNormalVerify == constant.SwitchOn {
				needVerify = true
			}
		}
	}

	var count int64
	if count, err = friendRepo.FriendRepo.Count(userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db count error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if params.ContactsFriendLimit > 0 && count >= params.ContactsFriendLimit {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("count max, count: %d, limit count: %d", count, params.ContactsFriendLimit))
		err = response.GetError(response.ErrFriendIsMax, lang)
		return
	}
	return
}

func (s *permissionUseCase) CheckDeleteFriendPermission(operationID string, userID string, lang string) (err error) {
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

		if params.IsMemberDelFriend == constant.SwitchOff {
			err = response.GetError(response.ErrFriendNormalNotDel, lang)
			return
		}
	}
	return
}
