package repo

import (
	"context"
	"im/pkg/db"
	"time"
)

var SettingCache = new(settingCache)

type settingCache struct{}

func (d *settingCache) IsAccountCodeExists(account string) (bool, error) {
	key := db.AccountTempCode + account
	n, err := db.RedisCli.Exists(context.Background(), key).Result()
	if n > 0 {
		return true, err
	}
	return false, err
}

func (d *settingCache) SetAccountCode(account string, code, ttl int) (err error) {
	key := db.AccountTempCode + account
	return db.RedisCli.Set(context.Background(), key, code, time.Duration(ttl)*time.Second).Err()
}

func (d *settingCache) GetAccountCode(account string) (string, error) {
	key := db.AccountTempCode + account
	return db.RedisCli.Get(context.Background(), key).Result()
}

func (d *settingCache) DelAccountCode(account string) (int64, error) {
	key := db.AccountTempCode + account
	return db.RedisCli.Del(context.Background(), key).Result()
}
