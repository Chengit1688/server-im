package usecase

import (
	"errors"
	"fmt"
	"im/internal/api/group/model"
	"im/internal/api/group/repo"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/go-redis/redis/v9"
)

var GroupMemberUseCase = new(groupMemberUseCase)

type groupMemberUseCase struct{}

func (c *groupMemberUseCase) GetMember(groupID string, memberID string) (member *model.GroupMember, err error) {
	member, err = repo.GroupMemberCache.GetMember(groupID, memberID)
	if err != nil && err != redis.Nil {
		return
	}

	if err == nil {
		return
	}

	key := repo.GroupMemberCache.GetMemberKey(groupID, memberID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	if member, err = repo.GroupMemberCache.GetMember(groupID, memberID); err == nil {
		return
	}

	if member, err = repo.GroupMemberRepo.GetMember(groupID, memberID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get member error, error: %v", err))
		return
	}

	if member.Id == 0 {
		err = errors.New("record not found")
		return
	}

	err = repo.GroupMemberCache.SetMember(groupID, member)
	return
}

func (c *groupMemberUseCase) SetMember(groupID string, memberID string) {
	repo.GroupMemberCache.DeleteMember(groupID, memberID)
}

func (c *groupMemberUseCase) CheckMember(groupID string, memberID string) bool {
	member, err := c.GetMember(groupID, memberID)
	if err != nil {
		return false
	}
	return member.Status == 1
}

func (c *groupMemberUseCase) UserInfoChange(userId string) (err error) {

	groupIdList, _ := db.CloumnList(model.GroupMember{}, model.GroupMember{
		UserId: userId,
		Status: 1,
	}, "group_id")

	for _, v := range groupIdList {
		groupId := v.(string)
		db.Update(model.GroupMember{}, model.GroupMember{
			GroupId: groupId,
			UserId:  userId,
			Status:  1,
		}, model.GroupMember{
			Version: 0,
		})
	}
	groupMemberList := []model.GroupMember{}
	total := int64(0)
	db.Find(model.GroupMember{}, model.GroupMember{
		UserId: userId,
		Status: 1,
	}, "", 1, 99999, &total, &groupMemberList)
	for _, v := range groupMemberList {
		groupId := v.GroupId

		c.SetMember(groupId, v.UserId)

		go repo.GroupRepo.GroupMemberChangeMsg("", groupId, []string{userId}, common.GroupMemberChangePush)
	}

	return
}

func (c *groupMemberUseCase) UserRoleChange(userId string, srcRole, optRole model.RoleType) (err error) {

	groupIdList, _ := db.CloumnList(model.GroupMember{}, model.GroupMember{
		UserId: userId,
		Status: 1,
		Role:   srcRole,
	}, "group_id")

	for _, v := range groupIdList {
		groupId := v.(string)
		db.Update(model.GroupMember{}, model.GroupMember{
			GroupId: groupId,
			UserId:  userId,
			Status:  1,
			Role:    srcRole,
		}, model.GroupMember{
			Role: optRole,
		})
	}

	groupMemberList := []model.GroupMember{}
	total := int64(0)
	db.Find(model.GroupMember{}, model.GroupMember{
		UserId: userId,
		Status: 1,
	}, "", 1, 99999, &total, &groupMemberList)
	for _, v := range groupMemberList {
		groupId := v.GroupId

		c.SetMember(groupId, v.UserId)

		go repo.GroupRepo.GroupMemberChangeMsg("", groupId, []string{userId}, common.GroupMemberChangePush)
	}

	return
}

func (c *groupMemberUseCase) SetMemberNickName(groupID, memberID, nickName string) error {
	return repo.GroupMemberRepo.SetMemberNickName(groupID, memberID, nickName)
}
