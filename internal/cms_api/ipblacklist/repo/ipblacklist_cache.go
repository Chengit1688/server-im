package repo

import (
	"context"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

const IpBlackListCacheKey = "cms_api:operation:ip_black_list"

var IpBlackListCache = new(ipblacklistCache)

type ipblacklistCache struct{}

func (c *ipblacklistCache) GetAllInfo() map[string]string {
	ctx := context.Background()
	value := db.RedisCli.HGetAll(ctx, IpBlackListCacheKey)

	infoMap, err := value.Result()
	if err != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", err.Error())
	}
	return infoMap
}

func (c *ipblacklistCache) Add(ip string) error {
	ctx := context.Background()
	value := db.RedisCli.HSet(ctx, IpBlackListCacheKey, ip, 1)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}

func (c *ipblacklistCache) Exist(ip string) (bool, error) {
	ctx := context.Background()
	value := db.RedisCli.HExists(ctx, IpBlackListCacheKey, ip)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Val(), value.Err()
}

func (c *ipblacklistCache) Del(ip string) error {
	ctx := context.Background()
	value := db.RedisCli.HDel(ctx, IpBlackListCacheKey, ip)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}
