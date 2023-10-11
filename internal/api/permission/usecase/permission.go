package usecase

import (
	"fmt"
	groupModel "im/internal/api/group/model"
	groupUseCase "im/internal/api/group/usecase"
	userModel "im/internal/api/user/model"
	userUseCase "im/internal/api/user/usecase"
	configModel "im/internal/cms_api/config/model"
	configUseCase "im/internal/cms_api/config/usecase"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
)

var PermissionUseCase = new(permissionUseCase)

type permissionUseCase struct{}

func (s *permissionUseCase) GetUserPermissionInfo(operationID string, userID string, lang string) (user *userModel.UserBaseInfo, params *configModel.ParameterConfigResp, err error) {
	if user, err = userUseCase.UserUseCase.GetBaseInfo(userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get user base info error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if params, err = configUseCase.ConfigUseCase.GetParameterConfig(); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}
	return
}

func (s *permissionUseCase) GetFriendPermissionInfo(operationID string, userID string, friendID string, lang string) (user *userModel.UserBaseInfo, friendUser *userModel.UserBaseInfo, params *configModel.ParameterConfigResp, err error) {
	if user, err = userUseCase.UserUseCase.GetBaseInfo(userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user get base info error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if friendUser, err = userUseCase.UserUseCase.GetBaseInfo(friendID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("friend get base info error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if params, err = configUseCase.ConfigUseCase.GetParameterConfig(); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}
	return
}

func (s *permissionUseCase) GetGroupAndMemberPermissionInfo(operationID string, userID string, groupID string, lang string) (group *groupModel.GroupInfo, member *groupModel.GroupMember, err error) {
	var groupInfo groupModel.GroupInfo
	if groupInfo, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group error, error: %v", err))
		err = response.GetError(response.ErrGroupNotExist, lang)
		return
	}

	if groupInfo.Status == 2 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group status = 2"))
		err = response.GetError(response.ErrGroupNotExist, lang)
		return
	}

	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	group = &groupInfo
	return
}

func (s *permissionUseCase) GetGroupMemberPermissionInfo(operationID string, userID string, groupID string, lang string) (member *groupModel.GroupMember, params *configModel.ParameterConfigResp, err error) {
	if member, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		err = response.GetError(response.ErrGroupNotMember, lang)
		return
	}

	if params, err = configUseCase.ConfigUseCase.GetParameterConfig(); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}
	return
}
