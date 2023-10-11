package repo

import (
	"errors"
	"fmt"
	chatModel "im/internal/api/chat/model"
	"im/internal/api/group/model"
	userModel "im/internal/api/user/model"
	userUse "im/internal/api/user/usecase"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	mqttMsg "im/pkg/mqtt"
	"im/pkg/util"
	"time"

	"gorm.io/gorm"
)

var GroupRepo = new(groupRepo)

type groupRepo struct{}

func (r *groupRepo) Search(keyword string, offset int, limit int) (list []model.Group, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)
	var listDB *gorm.DB
	listDB = db.DB.Model(&model.Group{}).Where("(group_id = ? OR name = ?) AND status = 1", keyword, keyword)

	if needPaging {
		listDB = listDB.Offset(offset).Limit(limit)
	}

	if err = listDB.Find(&list).Error; err != nil {
		return
	}

	if needPaging {
		if err = listDB.Count(&count).Error; err != nil {
			return
		}
	}
	return
}

func (r *groupRepo) Remove(groupID string, version int) (err error) {
	err = db.DB.Transaction(func(tx *gorm.DB) error {

		if err = tx.Model(&model.Group{}).Where("group_id = ?", groupID).Updates(map[string]interface{}{
			"status":       2,
			"last_version": version,
		}).Error; err != nil {
			return err
		}

		if err = tx.Model(&model.GroupMember{}).Where("group_id = ?", groupID).Updates(map[string]interface{}{
			"status": 2,
		}).Error; err != nil {
			return err
		}

		if err = tx.Table("conversations").Where("conversation_type = ? AND conversation_id = ?", chatModel.ConversationTypeGroup, groupID).Updates(map[string]interface{}{
			"deleted_at": time.Now().Unix(),
			"version":    -1,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

func (r *groupRepo) GetGroupRole(groups *[]model.GroupInfo, userId string) {
	for k, v := range *groups {
		memberInfo := model.GroupMember{}
		db.Info(&memberInfo, model.GroupMember{
			GroupId: v.GroupId,
			UserId:  userId,
			Status:  1,
		})
		(*groups)[k].Role = memberInfo.Role
		(*groups)[k].GroupNickName = memberInfo.GroupNickName
	}
}

func (r *groupRepo) GetGroupRoleOne(group *model.GroupInfo, userId string) {

	memberInfo := model.GroupMember{}
	db.Info(&memberInfo, model.GroupMember{
		GroupId: group.GroupId,
		UserId:  userId,
		Status:  1,
	})
	group.Role = memberInfo.Role
	group.GroupNickName = memberInfo.GroupNickName

}

func (r *groupRepo) GetGroupMemberUserInfo(members *[]model.GroupMemberInfo) {
	for k, v := range *members {

		userInfo, err := userUse.UserUseCase.GetBaseInfo(v.UserId)
		if err != nil {
			continue
		}
		(*members)[k].NickName = userInfo.NickName
		(*members)[k].FaceUrl = userInfo.FaceURL
		(*members)[k].BigFaceUrl = userInfo.BigFaceURL
		(*members)[k].Account = userInfo.Account
	}
}

func (r *groupRepo) GetApplyUserInfo(members *[]model.ApplyInfo) {
	for k, v := range *members {

		userInfo, err := userUse.UserUseCase.GetBaseInfo(v.UserId)
		if err == nil {
			(*members)[k].NickName = userInfo.NickName
			(*members)[k].FaceUrl = userInfo.FaceURL
			(*members)[k].BigFaceUrl = userInfo.BigFaceURL
		}

		groupInfo, _ := GroupCache.GroupInfo(v.GroupId)
		(*members)[k].GroupName = groupInfo.Name
	}
}

func (r *groupRepo) GetUserIdListByNickName(name string) (res []string) {
	uidList, err := db.CloumnList(userModel.User{}, map[string]interface{}{
		"nick_name": "?" + name,
	}, "user_id")
	if err != nil {
		return res
	}
	for _, v := range uidList {
		res = append(res, v.(string))
	}
	return res
}

func (r *groupRepo) GroupChangeMsg(operationID string, groupId string, code int) error {
	groupInfo := model.Group{}
	err := db.Info(&groupInfo, groupId)
	if err != nil {
		return err
	}
	groupInfoOut := model.GroupInfo{}
	util.CopyStructFields(&groupInfoOut, &groupInfo)
	msgCode := common.MessageType(0)
	switch code {
	case 1:
		msgCode = common.GroupAddPush
	case 2:
		msgCode = common.GroupRemovePush
	case 3:
		msgCode = common.GroupChangePush
	default:
		return errors.New("错误的类型")
	}
	startTime := time.Now().UnixMilli()
	mqttMsg.SendMessageToGroups(operationID, msgCode, groupInfoOut, groupId)
	endTime := time.Now().UnixMilli()
	logger.Sugar.Debug(operationID, "通知用时统计", fmt.Sprintf("群变化通知用时：%d", endTime-startTime))

	return nil
}

func (r *groupRepo) GroupChangeToUsersMsg(operationID string, groupId string, code int, userIdList []string) error {
	groupInfo := model.Group{}
	err := db.Info(&groupInfo, groupId)
	if err != nil {
		return err
	}
	groupInfoOut := model.GroupInfo{}
	util.CopyStructFields(&groupInfoOut, &groupInfo)
	msgCode := common.MessageType(0)
	switch code {
	case 1:
		msgCode = common.GroupAddPush
	case 2:
		msgCode = common.GroupRemovePush
	case 3:
		msgCode = common.GroupChangePush
	default:
		return errors.New("错误的类型")
	}
	startTime := time.Now().UnixMilli()
	mqttMsg.SendMessageToUsers(operationID, msgCode, groupInfoOut, userIdList...)
	endTime := time.Now().UnixMilli()
	if endTime-startTime > 1000 {
		logger.Sugar.Debug(operationID, "通知用时统计", fmt.Sprintf("群变化通知用时：%d", endTime-startTime))
	}

	return nil
}

func (r *groupRepo) GroupMemberChangeMsg(operationID string, groupId string, userIdList []string, mt common.MessageType) error {
	members := []model.GroupMember{}
	total := int64(0)
	err := db.Find(model.GroupMember{}, map[string]interface{}{
		"group_id": groupId,
		"user_id":  userIdList,
	}, "version asc", 1, 9999, &total, &members)
	if err != nil {
		return errors.New("获取变化用户失败")
	}
	memberInfos := []model.GroupMemberInfo{}
	for _, v := range members {
		temp := model.GroupMemberInfo{}
		util.CopyStructFields(&temp, v)
		memberInfos = append(memberInfos, temp)
	}
	r.GetGroupMemberUserInfo(&memberInfos)
	startTime := time.Now().UnixMilli()
	mqttMsg.SendMessageToGroups(operationID, mt, memberInfos, groupId)
	endTime := time.Now().UnixMilli()
	if endTime-startTime > 1000 {
		logger.Sugar.Debug(operationID, "通知用时统计", fmt.Sprintf("群成员变化通知用时：%d", endTime-startTime))
	}
	return nil
}

func (r *groupRepo) GroupMemberApplyChangeMsg(operationID string, applyId int64, code int) error {
	applyInfo := model.GroupMemberApply{}
	err := db.Info(&applyInfo, applyId)
	if err != nil {
		return errors.New("未查找到请求数据")
	}
	applyOut := model.ApplyInfo{}
	util.CopyStructFields(&applyOut, &applyInfo)
	applys := []model.ApplyInfo{
		applyOut,
	}
	r.GetApplyUserInfo(&applys)
	msgCode := common.MessageType(0)
	switch code {
	case 1:
		msgCode = common.GroupMemberApplyPush
	case 2:
		msgCode = common.GroupMemberVerifyPush
	default:
		return errors.New("错误的类型")
	}
	uidList := r.GroupAdminIdList(applyInfo.GroupId)
	startTime := time.Now().UnixMilli()
	mqttMsg.SendMessageToUsers(operationID, msgCode, applys[0], uidList...)
	logger.Sugar.Debugw(operationID, fmt.Sprintf("code:%d 发送群申请消息：%+v, 用户:%+v", msgCode, applys[0], uidList))
	endTime := time.Now().UnixMilli()
	if endTime-startTime > 1000 {
		logger.Sugar.Debug(operationID, "通知用时统计", fmt.Sprintf("群申请变化通知用时：%d", endTime-startTime))
	}
	return nil
}

func (r *groupRepo) GroupMemberIdList(groupId string) (idList []string) {
	return GroupCache.GroupMemberIdList(groupId)
}

func (r *groupRepo) GroupAdminIdList(groupId string) (idList []string) {
	idList = []string{}
	memberIdList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"group_id": groupId,
		"status":   1,
		"role":     []model.RoleType{model.RoleTypeAdmin, model.RoleTypeOwner},
	}, "user_id")
	if err != nil {
		return
	}
	for _, v := range memberIdList {
		idList = append(idList, v.(string))
	}
	return
}

func (r *groupRepo) UpdateMuteInfo(req model.GroupMuteAllReq) (err error) {
	updateMap := map[string]interface{}{}
	updateMap["mute_all_member"] = req.MuteAllMember
	updateMap["mute_all_period"] = req.MuteAllPeriod

	if err = db.DB.Model(&model.Group{}).Where("group_id = ?", req.GroupId).Updates(updateMap).Error; err != nil {
		return err
	}
	return nil
}
