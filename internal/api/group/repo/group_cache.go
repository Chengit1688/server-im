package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"im/internal/api/group/model"
	"im/pkg/db"
	"strings"
	"time"
)

var GroupCache = new(groupCache)

type groupCache struct{}

var groupVersionChan = make(chan bool, 1)

type groupVersionType struct {
	MemberVersion int
	GroupVersion  int
}

var groupVersionMap = map[string]*groupVersionType{}

func (c *groupCache) GroupMemberIdList(groupId string) (idList []string) {
	key := fmt.Sprintf("group_member_id_list:%s:cache", groupId)
	value := db.RedisCli.Get(context.Background(), key)
	if value.Err() != nil || value.Val() == "" {
		idList = c.UpGroupMemberIdCache(groupId)
		return idList
	}

	idList = strings.Split(value.Val(), ",")

	return
}

func (c *groupCache) UpGroupMemberIdCache(groupId string) []string {
	key := fmt.Sprintf("group_member_id_list:%s:cache", groupId)
	idList := []string{}
	memberIdList, err := db.CloumnList(&model.GroupMember{}, model.GroupMember{
		GroupId: groupId,
		Status:  1,
	}, "user_id")
	if err != nil {
		db.RedisCli.Del(context.Background(), key)
		return idList
	}
	for _, v := range memberIdList {
		idList = append(idList, v.(string))
	}
	db.RedisCli.Set(context.Background(), key, strings.Join(idList, ","), 3600*time.Second)
	return idList
}

func (c *groupCache) UpGroupInfoCache(groupId string) (res model.Group, err error) {
	key := fmt.Sprintf("group_info:%s:cache", groupId)
	err = db.Info(&res, groupId)
	if err != nil {
		return
	}
	infoJson, err := json.Marshal(res)
	if err != nil {
		return
	}
	db.RedisCli.Set(context.Background(), key, string(infoJson), 3600*time.Second)
	return
}

func (c *groupCache) GroupInfo(groupId string) (res model.Group, err error) {
	key := fmt.Sprintf("group_info:%s:cache", groupId)
	value := db.RedisCli.Get(context.Background(), key)
	if value.Err() != nil || value.Val() == "" {
		res, err = c.UpGroupInfoCache(groupId)
		return
	}

	err = json.Unmarshal([]byte(value.Val()), &res)
	return
}
