package repo

import (
	"im/internal/api/group/model"
	"im/pkg/db"
)

var GroupMemberRepo = new(groupMemberRepo)

type groupMemberRepo struct{}

func (r *groupMemberRepo) GroupCount(userID string) (count int64, err error) {
	err = db.DB.Model(&model.GroupMember{}).Where("user_id = ? AND role = ? AND status = 1", userID, model.RoleTypeOwner).Count(&count).Error
	return
}

func (r *groupMemberRepo) GetMember(groupID string, memberID string) (member *model.GroupMember, err error) {
	member = new(model.GroupMember)
	err = db.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ? AND status = 1", groupID, memberID).Find(member).Error
	return
}

func (r *groupMemberRepo) SetMemberNickName(groupID, memberID, nickName string) error {
	info := map[string]interface{}{
		"group_nick_name": nickName,
	}
	return db.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ? AND status = 1", groupID, memberID).Updates(info).Error
}

func (r *groupMemberRepo) GetAdminOwner(groupID string) (member []model.GroupMember, err error) {
	err = db.DB.Model(&model.GroupMember{}).Where("group_id = ? AND role IN ? AND status = 1", groupID, []string{string(model.RoleTypeAdmin), string(model.RoleTypeOwner)}).Find(&member).Error
	return
}

func (r *groupMemberRepo) UpdateMember(groupID, userID string, data interface{}) (err error) {
	err = db.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Updates(&data).Error
	return
}
