package repo

import (
	"context"
	"errors"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/go-redis/redis/v8"
)

const IpWhiteListCacheKey = "cms_api:operation:ip_white_list"

const IpWhiteListIsOpenCacheKey = "cms_api:operation:ip_white_list:is_open"
const IpWhiteListRootIDCacheKey = "cms_api:operation:ip_white_list:root_id"

var IpWhiteListCache = new(ipWhitelistCache)

type ipWhitelistCache struct{}

func (c *ipWhitelistCache) GetAllInfo() map[string]string {
	ctx := context.Background()
	value := db.RedisCli.HGetAll(ctx, IpWhiteListCacheKey)

	infoMap, err := value.Result()
	if err != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", err.Error())
	}
	return infoMap
}

func (c *ipWhitelistCache) Add(ip string) error {
	ctx := context.Background()
	value := db.RedisCli.HSet(ctx, IpWhiteListCacheKey, ip, 1)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}

func (c *ipWhitelistCache) Exist(ip string) (bool, error) {
	ctx := context.Background()
	value := db.RedisCli.HExists(ctx, IpWhiteListCacheKey, ip)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Val(), value.Err()
}

func (c *ipWhitelistCache) Del(ip string) error {
	ctx := context.Background()
	value := db.RedisCli.HDel(ctx, IpWhiteListCacheKey, ip)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}

func (c *ipWhitelistCache) Open() error {
	ctx := context.Background()
	value := db.RedisCli.Set(ctx, IpWhiteListIsOpenCacheKey, 1, 0)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}

func (c *ipWhitelistCache) Close() error {
	ctx := context.Background()
	value := db.RedisCli.Set(ctx, IpWhiteListIsOpenCacheKey, 2, 0)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}

func (c *ipWhitelistCache) GetIsOpen() int {
	ctx := context.Background()
	value, err := db.RedisCli.Get(ctx, IpWhiteListIsOpenCacheKey).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.Close()
			return 2
		}
	}
	return value
}

func (c *ipWhitelistCache) GetRootID() string {
	ctx := context.Background()
	value, err := db.RedisCli.Get(ctx, IpWhiteListRootIDCacheKey).Result()
	if err != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", err)
	}
	return value
}

func (c *ipWhitelistCache) SetRootID(rootID string) error {
	ctx := context.Background()
	value := db.RedisCli.Set(ctx, IpWhiteListRootIDCacheKey, rootID, 0)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}
