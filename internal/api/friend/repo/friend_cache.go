package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"im/internal/api/friend/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

var FriendCache = new(friendCache)
var friendVersionChan = make(chan bool, 1)

type friendCache struct{}

func (c *friendCache) GetFriendKey(userID string, friendID string) string {
	return fmt.Sprintf("%s:%s:%s", db.FriendInfo, userID, friendID)
}

func (c *friendCache) GetFriend(userID string, friendID string) (friend *model.Friend, err error) {
	key := c.GetFriendKey(userID, friendID)

	data, err := db.RedisCli.Get(context.Background(), key).Result()
	if err != nil {
		return
	}

	friend = new(model.Friend)
	err = util.JsonUnmarshal([]byte(data), friend)
	return
}

func (c *friendCache) SetFriend(userID string, friend *model.Friend) (err error) {
	key := c.GetFriendKey(userID, friend.FriendUserID)

	data, err := json.Marshal(friend)
	if err != nil {
		return
	}

	err = db.RedisCli.Set(context.Background(), key, string(data), util.RandDuration(util.OneWeek)).Err()
	return
}

func (c *friendCache) DeleteFriend(userID string, friendID string) {
	key := c.GetFriendKey(userID, friendID)

	if err := db.RedisCli.Del(context.Background(), key).Err(); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("redis del error, error: %v", err))
	}
}

func (c *friendCache) GetFriendMaxVersion(userId string) int {
	friendVersionChan <- true
	back := context.Background()
	key := fmt.Sprintf("friendversion:[%s]", userId)
	version, _ := db.RedisCli.Get(back, key).Int()
	if version < 10 {

		total := int64(0)
		friends := []model.Friend{}
		wheres := map[string]interface{}{
			"owner_user_id": userId,
		}
		err := db.Find(model.Friend{}, wheres, "version desc", 1, 1, &total, &friends)
		if err != nil {

			db.RedisCli.IncrBy(back, key, 10)
		} else {
			if len(friends) == 1 {
				if friends[0].Version < 10 {
					friends[0].Version = 10
				}
				db.RedisCli.IncrBy(back, key, int64(friends[0].Version))
			}
		}
		version, _ = db.RedisCli.Get(back, key).Int()
	}

	<-friendVersionChan
	return version
}
