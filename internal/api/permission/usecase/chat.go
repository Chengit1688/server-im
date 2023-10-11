package usecase

import (
	"fmt"
	chatModel "im/internal/api/chat/model"
	groupModel "im/internal/api/group/model"
	groupUseCase "im/internal/api/group/usecase"
	userModel "im/internal/api/user/model"
	configModel "im/internal/cms_api/config/model"
	configUseCase "im/internal/cms_api/config/usecase"
	"im/pkg/common/constant"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
	"time"
)

func (s *permissionUseCase) CheckChatGroupPermission(operationID string, groupID string, userID string, lang string) (err error) {

	var group groupModel.GroupInfo
	if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrGroupNotExist, lang)
	}

	if group.Status == 2 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrGroupNotExist, lang)
	}

	var memberInfo *groupModel.GroupMember
	memberInfo, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, userID)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group member error, error: %v", err))
		return response.GetError(response.ErrGroupNotMember, lang)
	}

	switch memberInfo.Role {
	case groupModel.RoleTypeOwner, groupModel.RoleTypeAdmin, groupModel.RoleTypeStaff:

	default:

		if group.MuteAllMember == constant.MuteMemberOpen {
			return response.GetError(response.ErrGroupMuteAll, lang)
		}

		if group.MuteAllMember == constant.MuteMemberPeriod {
			if s.isInTimePeriod(group.MuteAllPeriod) {
				return response.GetError(response.ErrGroupMuteAllPeriod, lang)
			}
		}
	}

	if memberInfo.MuteEndTime > time.Now().Unix() {
		return response.GetError(response.ErrGroupMuteUser, lang)
	}
	return
}

func (s *permissionUseCase) isInTimePeriod(timeStr string) bool {

	times := strings.Split(timeStr, "-")
	if len(times) != 2 {
		return false
	}
	startTimeStr := times[0]
	endTimeStr := times[1]

	if startTimeStr == endTimeStr {
		return true
	}

	dateStr := time.Now().Format("2006-01-02")
	location, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006-01-02T15:04:05", dateStr+"T"+startTimeStr+":00", location)
	endTime, _ := time.ParseInLocation("2006-01-02T15:04:05", dateStr+"T"+endTimeStr+":00", location)

	if endTime.Before(startTime) {

		endTime = endTime.Add(24 * time.Hour)
	}

	currentTime := time.Now()
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
}

func (s *permissionUseCase) CheckChatMultiSendPermission(operationID string, userID string, lang string) (err error) {
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

		if params.IsNormalMulSend == constant.SwitchOff {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("multi send error"))
			return response.GetError(response.ErrUserPermissions, lang)
		}
	}
	return
}

func (s *permissionUseCase) CheckChatRevokePermission(operationID string, userID string, sendID string, sendTime int64, lang string) (err error) {
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

		if userID != sendID {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("message sender not self"))
			err = response.GetError(response.ErrChatRevoke, lang)
			return
		}

		if params.RevokeTime != 0 && util.UnixMilliTime(time.Now())-sendTime > params.RevokeTime*1000 {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("revoke timeout"))
			return response.GetError(response.ErrChatRevoke, lang)
		}
	}
	return
}

func (s *permissionUseCase) CheckChatGroupRevokePermission(operationID string, groupID string, userID string, sendID string, sendTime int64, lang string) (err error) {
	var (
		member *groupModel.GroupMember
		params *configModel.ParameterConfigResp
	)

	var group groupModel.GroupInfo
	if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrGroupNotExist, lang)
	}

	if group.Status == 2 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrGroupNotExist, lang)
	}

	if member, params, err = s.GetGroupMemberPermissionInfo(operationID, userID, groupID, lang); err != nil {
		return
	}

	switch member.Role {
	case groupModel.RoleTypeOwner, groupModel.RoleTypeAdmin:

		return

	case groupModel.RoleTypeStaff:

		if userID == sendID {
			return
		}

		var sender *groupModel.GroupMember
		if sender, err = groupUseCase.GroupMemberUseCase.GetMember(groupID, sendID); err != nil || sender.Role == groupModel.RoleTypeUser {
			return
		}

		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("role: %s, sender role: %s", member.Role, sender.Role))

	case groupModel.RoleTypeUser:

		if userID == sendID && ((params.RevokeTime != 0 && util.UnixMilliTime(time.Now())-sendTime <= params.RevokeTime*1000) || params.RevokeTime == 0) {
			return
		}

		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("role: %s, revoke time: %ds", member.Role, params.RevokeTime))
	}

	err = response.GetError(response.ErrChatRevoke, lang)
	return
}

func (s *permissionUseCase) CheckChatRTCPermission(operationID string, rtcType chatModel.RTCType, lang string) (err error) {
	var params *configModel.ParameterConfigResp
	if params, err = configUseCase.ConfigUseCase.GetParameterConfig(); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
		err = response.GetError(response.ErrDB, lang)
		return
	}

	if rtcType == chatModel.RTCTypeAudio && params.IsOpenVoiceCall == constant.SwitchOff {
		err = response.GetError(response.ErrUserPermissions, lang)
		return
	}

	if rtcType == chatModel.RTCTypeVideo && params.IsOpenCameraCall == constant.SwitchOff {
		err = response.GetError(response.ErrUserPermissions, lang)
		return
	}
	return
}
