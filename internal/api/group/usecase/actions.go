package usecase

import (
	"errors"
	"fmt"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	chatUseCase "im/internal/api/chat/usecase"
	"im/internal/api/group/model"
	"im/internal/api/group/repo"
	userUse "im/internal/api/user/usecase"
	configModel "im/internal/cms_api/config/model"
	configUseCase "im/internal/cms_api/config/usecase"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
	"time"
)

func (c *groupUseCase) JoinGroup(operationID string, groupId, userId string) error {
	if groupId == "" || userId == "" {
		return errors.New("群id或用户id不能为空")
	}

	muteTime := int64(0)
	if groupInfo, err1 := repo.GroupCache.GroupInfo(groupId); err1 == nil {
		if groupInfo.Status == 2 {

			return nil
		}
		if groupInfo.MuteAllMember == constant.MuteNewJoinMember {
			muteTime = 365*24*3600 + time.Now().Unix()
		}
	}

	roleType := model.RoleTypeUser
	if userInfo, err := userUse.UserUseCase.GetBaseInfo(userId); err == nil {
		if userInfo.IsPrivilege == constant.SwitchOn {
			roleType = model.RoleTypeStaff
		}
	}

	err := db.Info(&model.GroupMember{}, &model.GroupMember{
		GroupId: groupId,
		UserId:  userId,
		Status:  1,
	})
	if err == nil {

		return nil
	}

	tx := db.DB.Begin()
	err = db.DeleteTx(tx, model.GroupMember{}, &model.GroupMember{
		GroupId: groupId,
		UserId:  userId,
	})
	if err != nil {
		return err
	}

	err = db.InsertTx(tx, &model.GroupMember{
		GroupId:      groupId,
		UserId:       userId,
		Status:       1,
		CreateTime:   time.Now().Unix(),
		Role:         roleType,
		JoinType:     "join",
		JoinByUserId: "",
		MuteEndTime:  muteTime,
	})
	if err != nil {
		return err
	}

	tx.Commit()
	repo.GroupCache.UpGroupInfoCache(groupId)
	repo.GroupCache.UpGroupMemberIdCache(groupId)

	repo.GroupRepo.GroupChangeToUsersMsg(operationID, groupId, 1, []string{userId})
	repo.GroupRepo.GroupChangeMsg(operationID, groupId, 3)
	repo.GroupRepo.GroupMemberChangeMsg(operationID, groupId, []string{userId}, common.GroupMemberAddPush)

	var (
		message *chatModel.MessageInfo
		content chatModel.MessageContent
	)
	content.OperatorID = ""
	content.BeOperatorList = append(content.BeOperatorList, chatModel.MessageBeOperator{BeOperatorID: userId})

	if message, err = chatUseCase.MessageUseCase.SendSystemMessageToGroupAndGroupMembers(operationID, groupId, chatModel.MessageGroupAddMemberNotify, &content, userId); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send system message to group and group members error, error: %v", err))
		return err
	}
	var params *configModel.ParameterConfigResp

	if params, err = configUseCase.ConfigUseCase.GetParameterConfig(); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
		err = code.ErrDB
		return nil
	}
	seq := message.Seq - 1
	if params.ShowGroupHistoryChatNum > 0 {
		if message.Seq-params.ShowGroupHistoryChatNum > 0 {
			seq = message.Seq - params.ShowGroupHistoryChatNum
		} else {
			seq = 0
		}
	}

	if err = chatUseCase.ConversationUseCase.UpsertUsersStartSeq(operationID, chatModel.ConversationTypeGroup, groupId, seq, userId); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("upsert users start seq error, error: %v", err))
	}
	return nil
}

func (c *groupUseCase) BatchJoinGroup(operationID string, groupId string, userIdList []string, opUserId string) (members []model.GroupMember, err error) {
	hadUserList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"group_id": groupId,
		"user_id":  userIdList,
		"status":   1,
	}, "user_id")
	hadUserIdMap := map[string]bool{}
	if err == nil {
		for _, v := range hadUserList {
			uid := v.(string)
			hadUserIdMap[uid] = true
		}
	}
	newUidList := []string{}
	for _, v := range userIdList {
		if had := hadUserIdMap[v]; !had {
			newUidList = append(newUidList, v)
		}
	}
	if len(newUidList) == 0 {
		return
	}

	muteTime := int64(0)
	if groupInfo, err1 := repo.GroupCache.GroupInfo(groupId); err1 == nil {
		if groupInfo.MuteAllMember == constant.MuteNewJoinMember {
			muteTime = 365*24*3600 + time.Now().Unix()
		}
	}

	logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "筛选有效用户完毕")

	newMembers := []model.GroupMember{}
	nows := time.Now().Unix()
	for _, v := range newUidList {
		roleType := model.RoleTypeUser
		if userInfo, err := userUse.UserUseCase.GetBaseInfo(v); err == nil {
			if userInfo.IsPrivilege == constant.SwitchOn {
				roleType = model.RoleTypeStaff
			}
		}

		newMembers = append(newMembers, model.GroupMember{
			GroupId:      groupId,
			UserId:       v,
			Role:         roleType,
			RoleIndex:    300,
			JoinType:     "join",
			JoinByUserId: "",
			CreateTime:   nows,
			MuteEndTime:  muteTime,
		})
	}
	logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "开始邀请事务")
	tx := db.DB.Begin()
	err = db.DeleteTx(tx, model.GroupMember{}, map[string]interface{}{
		"group_id": groupId,
		"user_id":  newUidList,
	})
	if err != nil {
		return
	}
	err = db.InsertTx(tx, &newMembers)
	if err != nil {
		return
	}

	logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "结束邀请事务")
	tx.Commit()

	repo.GroupCache.UpGroupInfoCache(groupId)
	repo.GroupCache.UpGroupMemberIdCache(groupId)
	logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "清理缓存完毕")
	go func() {

		logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "开始发送消息1")
		repo.GroupRepo.GroupChangeToUsersMsg(operationID, groupId, 1, newUidList)
		logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "开始发送消息2")
		repo.GroupRepo.GroupChangeMsg(operationID, groupId, 3)
		logger.Sugar.Debug(operationID, util.GetSelfFuncName(), "开始发送消息3")
		repo.GroupRepo.GroupMemberChangeMsg(operationID, groupId, newUidList, common.GroupMemberAddPush)

		var (
			message *chatModel.MessageInfo
			content chatModel.MessageContent
		)
		content.OperatorID = opUserId
		for _, member := range newMembers {
			content.BeOperatorList = append(content.BeOperatorList, chatModel.MessageBeOperator{BeOperatorID: member.UserId})
		}

		if message, err = chatUseCase.MessageUseCase.SendSystemMessageToGroupAndGroupMembers(operationID, groupId, chatModel.MessageGroupAddMemberNotify, &content, newUidList...); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send system message to group and group members error, error: %v", err))
			return
		}
		var params *configModel.ParameterConfigResp

		if params, err = configUseCase.ConfigUseCase.GetParameterConfig(); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get parameter config error, error: %v", err))
			err = code.ErrDB
			return
		}
		seq := message.Seq - 1
		if params.ShowGroupHistoryChatNum > 0 {
			if message.Seq-params.ShowGroupHistoryChatNum > 0 {
				seq = message.Seq - params.ShowGroupHistoryChatNum
			} else {
				seq = 0
			}
		}

		if err = chatUseCase.ConversationUseCase.UpsertUsersStartSeq(operationID, chatModel.ConversationTypeGroup, groupId, seq, newUidList...); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("upsert users start seq error, error: %v", err))
		}
	}()
	return newMembers, nil
}

func (c *groupUseCase) CreateGroup(operationID, name, faceurl, userId string, isTop int, notifyCation ...string) (group model.Group, err error) {
	tx := db.DB.Begin()
	groupId := util.RandID(db.GroupIDSize)
	newGroup := model.Group{
		GroupId:            groupId,
		Name:               name,
		CreateUserId:       userId,
		FaceUrl:            faceurl,
		CreateTime:         time.Now().Unix(),
		LastVersion:        10,
		LastMemberVersion:  10,
		MembersTotal:       1,
		AdminsTotal:        1,
		IsTopannocuncement: isTop,
	}
	if len(notifyCation) > 0 {
		newGroup.Notification = notifyCation[0]
	}
	logger.Sugar.Debugw("随机群id", groupId, newGroup)
	err = db.InsertTx(tx, &newGroup)
	if err != nil {
		return
	}
	newData := model.GroupMember{
		GroupId:    groupId,
		UserId:     userId,
		Role:       model.RoleTypeOwner,
		RoleIndex:  100,
		CreateTime: time.Now().Unix(),
		Version:    10,
	}
	err = db.InsertTx(tx, &newData)
	if err != nil {
		return
	}

	tx.Commit()

	repo.GroupCache.UpGroupInfoCache(groupId)
	repo.GroupCache.UpGroupMemberIdCache(groupId)

	repo.GroupRepo.GroupChangeToUsersMsg(operationID, groupId, 1, []string{userId})

	repo.GroupRepo.GroupMemberChangeMsg(operationID, groupId, []string{userId}, common.GroupMemberAddPush)

	var (
		message *chatModel.MessageInfo
		content chatModel.MessageContent
	)
	content.OperatorID = userId

	if message, err = chatUseCase.MessageUseCase.SendSystemMessageToUsers(operationID, chatModel.ConversationTypeGroup, newData.GroupId, chatModel.MessageGroupCreateNotify, &content, userId); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send system message to users error, error: %v", err))
		return
	}

	if err = chatUseCase.ConversationUseCase.UpsertUsersStartSeq(operationID, chatModel.ConversationTypeGroup, newData.GroupId, message.Seq-1, userId); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("upsert users start seq error, error: %v", err))
	}
	return newGroup, nil
}

func (c *groupUseCase) UpdateGroup(operationID, groupId string, updateInfo *model.Group, updateMap map[string]interface{}, opUserID string) (group model.Group, err error) {

	logger.Sugar.Debugw("群更新", updateInfo)
	groupStatus := 0
	muteAllStatus := 0
	notify := ""
	if updateInfo != nil {
		groupStatus = updateInfo.Status
		muteAllStatus = updateInfo.MuteAllMember
		notify = updateInfo.Notification
	}

	db.Info(&group, groupId)
	isSendNotify := false
	if group.Notification != notify {
		isSendNotify = true
	}

	if updateMap != nil {
		if v, had := updateMap["status"]; had {
			groupStatus = v.(int)
		}
		if v, had := updateMap["mute_all_member"]; had {
			muteAllStatus = v.(int)
		}
		if v, had := updateMap["notification"]; had {
			notify = v.(string)
		}
	}

	if updateInfo == nil && updateMap == nil {
		return group, errors.New("更新参数为空")
	}

	updateNullInfo := map[string]interface{}{}
	if updateInfo != nil && updateInfo.Notification != "" && strings.TrimSpace(updateInfo.Notification) == "" {
		updateNullInfo["notification"] = ""
	}

	if updateInfo != nil && updateInfo.Introduction != "" && strings.TrimSpace(updateInfo.Introduction) == "" {
		updateNullInfo["introduction"] = ""
	}

	if updateInfo != nil && updateInfo.GroupSendLimit < 0 {
		updateNullInfo["group_send_limit"] = 0
	}

	if updateInfo != nil && updateInfo.IsDisplayNicknameOpen == 0 {
		updateNullInfo["is_display_nickname_open"] = 1
	}

	if updateInfo != nil && updateInfo.MuteAllMember != constant.MuteMemberPeriod {
		updateNullInfo["mute_all_period"] = ""
	}
	updateAt := time.Now().Unix()

	tx := db.DB.Begin()
	if updateInfo != nil {
		updateInfo.UpdatedAt = updateAt
		err = db.UpdateTx(tx, &model.Group{}, groupId, updateInfo)
		if err != nil {
			return
		}
	}

	if updateMap != nil {
		updateMap["updated_at"] = updateAt
		err = db.UpdateTx(tx, &model.Group{}, groupId, updateMap)
		if err != nil {
			return
		}
	}

	if len(updateNullInfo) != 0 {

		err = db.UpdateTx(tx, &model.Group{}, groupId, updateNullInfo)
		if err != nil {
			return
		}
	}

	if groupStatus == 2 {
		err = db.UpdateTx(tx, &model.Group{}, groupId, map[string]interface{}{
			"admins_total":  0,
			"members_total": 0,
			"updated_at":    updateAt,
		})
		if err != nil {
			return
		}
	}
	oldMemberIdList := repo.GroupRepo.GroupMemberIdList(groupId)
	if groupStatus == 2 {
		err = db.UpdateTx(tx, &model.GroupMember{}, model.GroupMember{GroupId: groupId, Status: 1}, model.GroupMember{
			Status: 2,
		})
		if err != nil {
			logger.Sugar.Error("取消所有群成员失败", err)
			return
		}

		err = db.UpdateTx(tx, &model.GroupMemberApply{}, map[string]interface{}{
			"status":   0,
			"group_id": updateInfo.GroupId,
		}, model.GroupMemberApply{
			Status: 2,
		})
		if err != nil {
			logger.Sugar.Error("取消所有群申请失败", err)
			return
		}

		if err2 := chatRepo.ConversationRepo.Delete(chatModel.ConversationTypeGroup, groupId); err2 != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("delete group conversation error, error: %v", err2))
		}
	}
	tx.Commit()

	db.Info(&group, groupId)
	if groupStatus == 2 {

		repo.GroupCache.UpGroupInfoCache(groupId)
		repo.GroupCache.UpGroupMemberIdCache(groupId)
		go repo.GroupRepo.GroupChangeToUsersMsg(operationID, groupId, 2, oldMemberIdList)

	} else {

		repo.GroupCache.UpGroupInfoCache(groupId)
		repo.GroupRepo.GroupChangeMsg(operationID, groupId, 3)
	}

	var content chatModel.MessageContent
	content.OperatorID = opUserID
	if muteAllStatus == 1 {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupAllMuteNotify, &content)
	}
	if muteAllStatus == 2 {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupAllUnmuteNotify, &content)
	}
	if strings.TrimSpace(notify) != "" && isSendNotify {
		var contentNotify chatModel.MessageContent
		contentNotify.OperatorID = opUserID
		contentNotify.OperatorFaceUrl = ""
		contentNotify.Content = notify
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupNotifyChangeNotify, &contentNotify)
	}

	return
}

func (c *groupUseCase) UpdateInformation(operationID string, userID string, groupID string, group *model.GroupInformationInfo) (newGroup model.Group, err error) {
	var _ map[string]interface{}
	if _, err = util.StructToMapWithoutNil(group, "json"); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("struct to map error, error: %v", err))
		err = code.ErrUnknown
		return
	}

	userInfo, err := userUse.UserUseCase.GetBaseInfo(userID)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetBaseInfo error, error: %v", err))
		err = code.ErrDB
		return
	}

	if _, err = repo.GroupCache.UpGroupInfoCache(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group cache update error, error: %v", err))
		err = code.ErrDB
		return
	}

	newGroup = c.GroupInfoPush(operationID, groupID, common.GroupChangePush)

	if group.Notification != nil && *group.Notification != "" {
		var content chatModel.MessageContent
		content.OperatorID = userID
		content.OperatorFaceUrl = userInfo.FaceURL
		content.Content = *group.Notification
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupID, chatModel.MessageGroupNotifyChangeNotify, &content)
	}
	return
}

func (c *groupUseCase) UpdateManage(operationID string, groupID string, group *model.GroupManageInfo) (newGroup model.Group, err error) {
	var _ map[string]interface{}
	if _, err = util.StructToMapWithoutNil(group, "json"); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("struct to map error, error: %v", err))
		err = code.ErrUnknown
		return
	}

	if _, err = repo.GroupCache.UpGroupInfoCache(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group cache update error, error: %v", err))
		err = code.ErrDB
		return
	}

	newGroup = c.GroupInfoPush(operationID, groupID, common.GroupChangePush)
	return
}

func (c *groupUseCase) Remove(operationID string, groupID string) (err error) {
	if err = repo.GroupRepo.Remove(groupID, 0); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group remove error, error: %v", err))
		err = code.ErrDB
		return
	}

	if _, err = repo.GroupCache.UpGroupInfoCache(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group cache update error, error: %v", err))
		err = code.ErrDB
		return
	}

	c.GroupInfoPush(operationID, groupID, common.GroupRemovePush)
	return
}

func (c *groupUseCase) GroupRemoveMember(operationID, groupId string, userIdList []string, opUserId string, reason model.ReasonType) error {
	groupMemberStatus := 2
	if reason == model.ReasonTypeKick {
		groupMemberStatus = 3
	}

	hadUserList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"group_id": groupId,
		"user_id":  userIdList,
		"status":   1,
	}, "user_id")
	hadUserIdMap := map[string]bool{}

	ownerInfo := model.GroupMember{}
	db.Info(&ownerInfo, &model.GroupMember{
		GroupId: groupId,
		Role:    model.RoleTypeOwner,
		Status:  1,
	})
	if err == nil {
		for _, v := range hadUserList {
			uid := v.(string)
			if uid == ownerInfo.UserId {
				continue
			}
			hadUserIdMap[uid] = true
		}
	}
	removeUidList := []string{}
	for _, v := range userIdList {
		if had := hadUserIdMap[v]; had {
			removeUidList = append(removeUidList, v)
		}
	}

	if len(removeUidList) == 0 {
		return nil
	}

	tx := db.DB.Begin()
	for _, v := range removeUidList {
		err = db.UpdateTx(tx, &model.GroupMember{}, model.GroupMember{
			GroupId: groupId,
			UserId:  v,
			Status:  1,
		}, model.GroupMember{
			Status: groupMemberStatus,
		})
		if err != nil {
			return err
		}

		GroupMemberUseCase.SetMember(groupId, v)
	}

	tx.Commit()

	var content chatModel.MessageContent
	content.OperatorID = opUserId
	content.BeOperatorList = []chatModel.MessageBeOperator{}
	for _, v := range removeUidList {
		content.BeOperatorList = append(content.BeOperatorList, chatModel.MessageBeOperator{
			BeOperatorID: v,
		})
	}
	chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupDeleteNotify, &content)

	if err2 := chatRepo.ConversationRepo.Delete(chatModel.ConversationTypeGroup, groupId, removeUidList...); err2 != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("delete group members conversation error, error: %v", err2))
	}

	if reason == model.ReasonTypeKick {
		if err2 := chatRepo.MessageRepo.UpdateSenderMessageStatus(chatModel.ConversationTypeGroup, groupId, chatModel.MessageStatusTypeDelete, removeUidList...); err2 != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update group members message status error, error: %v", err2))
		}
	}

	repo.GroupCache.UpGroupInfoCache(groupId)
	repo.GroupCache.UpGroupMemberIdCache(groupId)

	repo.GroupRepo.GroupChangeToUsersMsg(operationID, groupId, 2, removeUidList)
	repo.GroupRepo.GroupChangeMsg(operationID, groupId, 3)
	repo.GroupRepo.GroupMemberChangeMsg(operationID, groupId, removeUidList, common.GroupMemberRemovePush)
	return nil
}

func (c *groupUseCase) SetGroupAdmin(operationID, groupId, userId string, status int, opUserId string) (err error) {
	var role model.RoleType
	roleIndex := 0
	memberInfo := model.GroupMember{}
	err = db.Info(&memberInfo, model.GroupMember{
		GroupId: groupId,
		UserId:  userId,
		Status:  1,
	})
	if err != nil {
		return
	}
	if status == 1 {
		if memberInfo.Role == model.RoleTypeAdmin {
			return
		}
		role = model.RoleTypeAdmin
		roleIndex = 200
	}

	roleType := model.RoleTypeUser
	if userInfo, err := userUse.UserUseCase.GetBaseInfo(userId); err == nil {
		if userInfo.IsPrivilege == constant.SwitchOn {
			roleType = model.RoleTypeStaff
		}
	}

	if status == 2 {
		if memberInfo.Role == roleType {
			return
		}
		role = roleType
		roleIndex = 300
	}

	tx := db.DB.Begin()
	err = db.UpdateTx(tx, &model.GroupMember{}, memberInfo.Id, model.GroupMember{
		Role:      role,
		RoleIndex: roleIndex,
	})
	if err != nil {
		return
	}

	tx.Commit()

	GroupMemberUseCase.SetMember(groupId, userId)

	var content chatModel.MessageContent
	content.OperatorID = opUserId
	content.BeOperatorList = append(content.BeOperatorList, chatModel.MessageBeOperator{BeOperatorID: userId})
	if status == 1 {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupSetAdminNotify, &content)
	} else {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupUnsetAdminNotify, &content)
	}

	repo.GroupRepo.GroupMemberChangeMsg(operationID, groupId, []string{userId}, common.GroupMemberChangePush)
	return nil
}

func (c *groupUseCase) SetGroupOwner(operationID, groupId, userId string, opUserId string) (err error) {
	ownerMemberInfo := model.GroupMember{}
	err = db.Info(&ownerMemberInfo, model.GroupMember{
		GroupId: groupId,
		Role:    model.RoleTypeOwner,
		Status:  1,
	})
	if err != nil {
		return
	}

	roleType := model.RoleTypeUser
	if userInfo, err := userUse.UserUseCase.GetBaseInfo(ownerMemberInfo.UserId); err == nil {
		if userInfo.IsPrivilege == constant.SwitchOn {
			roleType = model.RoleTypeStaff
		}
	}

	tx := db.DB.Begin()
	err = db.UpdateTx(tx, &model.GroupMember{}, ownerMemberInfo.Id, model.GroupMember{
		Role:      roleType,
		RoleIndex: 300,
	})
	if err != nil {
		return
	}
	role := model.RoleTypeOwner
	err = db.UpdateTx(tx, &model.GroupMember{}, model.GroupMember{
		GroupId: groupId,
		UserId:  userId,
		Status:  1,
	}, map[string]interface{}{
		"role":          role,
		"role_index":    100,
		"mute_end_time": 0,
	})

	if err != nil {
		return
	}

	tx.Commit()

	GroupMemberUseCase.SetMember(groupId, userId)
	GroupMemberUseCase.SetMember(groupId, opUserId)

	var content chatModel.MessageContent
	content.OperatorID = opUserId
	content.BeOperatorList = append(content.BeOperatorList, chatModel.MessageBeOperator{BeOperatorID: userId})
	chatUseCase.MessageUseCase.SendSystemMessageToGroup(operationID, groupId, chatModel.MessageGroupTransferNotify, &content)

	repo.GroupRepo.GroupMemberChangeMsg(operationID, groupId, []string{userId, ownerMemberInfo.UserId}, common.GroupMemberChangePush)
	return nil
}

func (c *groupUseCase) UpdateGroupMuteInfo(req model.GroupMuteAllReq) (err error) {
	if err = repo.GroupRepo.UpdateMuteInfo(req); err != nil {
		return err
	}

	repo.GroupCache.UpGroupInfoCache(req.GroupId)
	repo.GroupRepo.GroupChangeMsg(req.OperationID, req.GroupId, 3)

	var content chatModel.MessageContent
	if req.MuteAllMember == constant.MuteMemberOpen || req.MuteAllMember == constant.MuteNewJoinMember || req.MuteAllMember == constant.MuteMemberPeriod {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(req.OperationID, req.GroupId, chatModel.MessageGroupAllMuteNotify, &content)
	}
	if req.MuteAllMember == constant.MuteMemberClose {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(req.OperationID, req.GroupId, chatModel.MessageGroupAllUnmuteNotify, &content)
	}

	return nil
}
