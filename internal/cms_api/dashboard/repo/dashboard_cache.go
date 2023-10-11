package repo

import (
	"context"
	"im/pkg/db"
)

var DashboardCache = new(dashboardCache)

type dashboardCache struct{}

func (c *dashboardCache) GetOnlineMax() int64 {
	key := "cms_api:dashboard:online_max"
	ctx := context.Background()
	value := db.RedisCli.Get(ctx, key)
	if value.Err() != nil {
		db.RedisCli.Set(ctx, key, 0, -1)
		return 0
	}
	num, _ := value.Int64()
	return num
}

func (c *dashboardCache) SetOnlineMax(num int64) error {
	key := "cms_api:dashboard:online_max"
	ctx := context.Background()
	value := db.RedisCli.Set(ctx, key, num, -1)
	return value.Err()
}
