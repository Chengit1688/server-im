package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"im/internal/api/group/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

var GroupMemberCache = new(groupMemberCache)

type groupMemberCache struct{}

func (c *groupMemberCache) GetMemberKey(groupID string, memberID string) string {
	return fmt.Sprintf("%s:%s:%s", db.GroupMemberInfo, groupID, memberID)
}

func (c *groupMemberCache) GetMember(groupID string, memberID string) (member *model.GroupMember, err error) {
	key := c.GetMemberKey(groupID, memberID)

	data, err := db.RedisCli.Get(context.Background(), key).Result()
	if err != nil {
		return
	}

	member = new(model.GroupMember)
	err = util.JsonUnmarshal([]byte(data), member)
	return
}

func (c *groupMemberCache) SetMember(groupID string, member *model.GroupMember) (err error) {
	key := c.GetMemberKey(groupID, member.UserId)

	data, err := json.Marshal(member)
	if err != nil {
		return
	}

	err = db.RedisCli.Set(context.Background(), key, string(data), util.RandDuration(util.OneWeek)).Err()
	return
}

func (c *groupMemberCache) DeleteMember(groupID string, memberID string) {
	key := c.GetMemberKey(groupID, memberID)

	if err := db.RedisCli.Del(context.Background(), key).Err(); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("redis del error, error: %v", err))
	}
}
