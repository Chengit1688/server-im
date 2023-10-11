package repo

import (
	"context"
	"im/pkg/db"
	"time"
)

var MultiSendLockKey = "cms_api:operation:multi_send"
var MessageCache = new(messageCache)

type messageCache struct{}

func (c *messageCache) GetMultiSendLock() (running bool, err error) {
	ctx := context.Background()
	value := db.RedisCli.Exists(ctx, MultiSendLockKey)
	running = value.Val() > 0
	err = value.Err()
	return
}

func (c *messageCache) SetMultiSendLock() (err error) {
	ctx := context.Background()
	value := db.RedisCli.SetEx(ctx, MultiSendLockKey, 1, time.Minute*10)
	err = value.Err()
	return
}

func (c *messageCache) DelMultiSendLock() (err error) {
	ctx := context.Background()
	value := db.RedisCli.Del(ctx, MultiSendLockKey)
	err = value.Err()
	return
}
