package usecase

import (
	"fmt"
	"im/internal/api/group/model"
	"im/internal/api/group/repo"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
)

var GroupUseCase = new(groupUseCase)

type groupUseCase struct{}

func (c *groupUseCase) Search(userID string, keyword string, offset int, limit int) (list []model.GroupInfo, count int64, err error) {
	var groups []model.Group
	if groups, count, err = repo.GroupRepo.Search(keyword, offset, limit); err != nil {
		return
	}

	for _, group := range groups {
		var data model.GroupInfo
		data.Group = group

		if member, err2 := GroupMemberUseCase.GetMember(group.GroupId, userID); err2 == nil {
			data.Role = member.Role
		}
		list = append(list, data)
	}
	return
}

func (c *groupUseCase) GroupInfo(groupId string) (info model.GroupInfo, err error) {

	group, err := repo.GroupCache.GroupInfo(groupId)
	if err == nil {
		util.CopyStructFields(&info, group)
	}
	return
}

func (c *groupUseCase) GroupMemberIdList(groupId string) (idList []string) {
	return repo.GroupRepo.GroupMemberIdList(groupId)
}

func (c *groupUseCase) GroupDelete(groupId string) error {
	tx := db.DB.Begin()
	err := db.DeleteTx(tx, model.Group{}, model.Group{
		GroupId: groupId,
	})
	if err != nil {
		return err
	}
	err = db.DeleteTx(tx, model.GroupMember{}, model.GroupMember{
		GroupId: groupId,
	})
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (c *groupUseCase) GetGroupVersions(userId string) (res []model.GroupsVersion) {
	res = []model.GroupsVersion{}
	groupIdList, err := db.CloumnList(model.GroupMember{}, model.GroupMember{
		UserId: userId,
		Status: 1,
	}, "group_id")
	if err != nil {
		return
	}
	groupIds := []string{}
	for _, v := range groupIdList {
		groupIds = append(groupIds, v.(string))
	}

	groups := []model.Group{}
	total := int64(0)
	err = db.Find(model.Group{}, map[string]interface{}{
		"group_id": groupIds,
		"status":   1,
	}, "", 1, 9999, &total, &groups)
	if err != nil {
		return
	}
	for _, v := range groups {
		res = append(res, model.GroupsVersion{
			GroupId:       v.GroupId,
			GroupVersion:  v.LastVersion,
			MemberVersion: v.LastMemberVersion,
		})
	}

	return
}

func (c *groupUseCase) GetDefaultGroups() (res []model.Group) {
	res = []model.Group{}
	total := int64(0)
	db.Find(model.Group{}, model.Group{
		IsDefault: 1,
		Status:    1,
	}, "", 1, 9999, &total, &res)
	return
}

func (c *groupUseCase) GroupInfoPush(operationID string, groupID string, mt common.MessageType, userIDList ...string) (group model.Group) {
	var (
		groupInfo model.GroupInfo
		err       error
	)
	if groupInfo, err = c.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get group info error, error: %v", err))
		return
	}

	if len(userIDList) != 0 {
		_ = mqtt.SendMessageToUsers(operationID, mt, groupInfo, userIDList...)
		return
	}

	_ = mqtt.SendMessageToGroups(operationID, mt, groupInfo, groupID)
	group = groupInfo.Group
	return
}
