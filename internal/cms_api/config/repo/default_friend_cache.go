package repo

import (
	"context"
	"fmt"
	"im/pkg/db"
)

var DefaultFriendCahce = new(defaultFriendCache)

type defaultFriendCache struct{}

func (c *defaultFriendCache) GetDefaultFriendKey() string {
	return fmt.Sprintf("%s", db.DefaultFriendVersion)
}

func (c *defaultFriendCache) GetVersion() (version int64, err error) {
	key := c.GetDefaultFriendKey()
	version, err = db.RedisCli.Get(context.Background(), key).Int64()
	if version <= 0 {
		err = db.RedisCli.SetNX(context.Background(), key, 0, 0).Err()
		return
	}
	return
}

func (c *defaultFriendCache) setVersion(value interface{}) (err error) {
	key := c.GetDefaultFriendKey()
	err = db.RedisCli.Set(context.Background(), key, value, 0).Err()
	return
}

func (c *defaultFriendCache) setVersionIncr() (version int64, err error) {
	key := c.GetDefaultFriendKey()
	_, err = db.RedisCli.Incr(context.Background(), key).Result()
	return
}
