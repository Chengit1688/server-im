package repo

import (
	"context"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

const DeviceListCacheKey = "cms_api:operation:device_block_list"

var DeviceListCache = new(devicelistCache)

type devicelistCache struct{}

func (c *devicelistCache) GetAllInfo() map[string]string {
	ctx := context.Background()
	value := db.RedisCli.HGetAll(ctx, DeviceListCacheKey)

	infoMap, err := value.Result()
	if err != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", err.Error())
	}
	return infoMap
}

func (c *devicelistCache) Add(deviceID string) error {
	ctx := context.Background()
	value := db.RedisCli.HSet(ctx, DeviceListCacheKey, deviceID, 1)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}

func (c *devicelistCache) Exist(deviceID string) (bool, error) {
	ctx := context.Background()
	value := db.RedisCli.HExists(ctx, DeviceListCacheKey, deviceID)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Val(), value.Err()
}

func (c *devicelistCache) Del(deviceID string) error {
	ctx := context.Background()
	value := db.RedisCli.HDel(ctx, DeviceListCacheKey, deviceID)
	if value.Err() != nil {
		logger.Sugar.Errorw("redis err", "func", util.GetSelfFuncName(), "error", value.Err())
	}
	return value.Err()
}
