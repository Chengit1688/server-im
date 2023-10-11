package repo

import (
	"fmt"
	groupModel "im/internal/api/group/model"
	userUse "im/internal/api/user/usecase"
	"im/internal/cms_api/group/model"
	"im/pkg/db"
)

var GroupRepo = new(groupRepo)

type groupRepo struct{}

func (r *groupRepo) GetGroupsOwnerInfo(groups *[]model.GroupInfo) {
	for k, v := range *groups {

		ownerIds, err := db.CloumnList(groupModel.GroupMember{}, groupModel.GroupMember{
			GroupId: v.GroupId,
			Role:    groupModel.RoleTypeOwner,
			Status:  1,
		}, "user_id")
		if err != nil || len(ownerIds) == 0 {
			continue
		}

		userInfo, err := userUse.UserUseCase.GetBaseInfo(ownerIds[0].(string))
		if err != nil {
			continue
		}
		(*groups)[k].OwnerNickName = userInfo.NickName
		(*groups)[k].OwnerUserId = ownerIds[0].(string)
	}
}

func (r *groupRepo) GetGroupOwnerInfo(group *model.GroupInfo) {

	ownerIds, err := db.CloumnList(groupModel.GroupMember{}, groupModel.GroupMember{
		GroupId: group.GroupId,
		Role:    groupModel.RoleTypeOwner,
		Status:  1,
	}, "user_id")
	if err != nil || len(ownerIds) == 0 {
		return
	}

	userInfo, err := userUse.UserUseCase.GetBaseInfo(ownerIds[0].(string))
	if err == nil {
		group.OwnerNickName = userInfo.NickName
		group.OwnerUserId = ownerIds[0].(string)
	}

}

func (r *groupRepo) GetGroupMemberUserInfo(members *[]model.GroupMemberInfo) {
	for k, v := range *members {

		userInfo, err := userUse.UserUseCase.GetBaseInfo(v.UserId)
		if err != nil {
			continue
		}
		(*members)[k].NickName = userInfo.NickName
		(*members)[k].FaceUrl = userInfo.FaceURL
		(*members)[k].Account = userInfo.Account
	}
}

func (r *groupRepo) GroupSearch(search string) (groups []groupModel.Group, count int64, err error) {
	tx := db.DB.Model(groupModel.Group{})
	tx = tx.Where("status = 1 AND (group_id LIKE ? OR name LIKE ?)", fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search))
	err = tx.Find(&groups).Count(&count).Error
	return
}
