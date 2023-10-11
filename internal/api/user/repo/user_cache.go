package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"im/config"
	momentsModel "im/internal/api/moments/model"
	"im/internal/api/user/model"
	cmsConfigModel "im/internal/cms_api/config/model"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"math"
	"time"
)

var UserCache = new(userCache)

type userCache struct{}

const (
	RetryReadRedisTimes = 2
)

func (u *userCache) AddTokenCache(userInfo *model.User) (string, error) {
	var (
		tokenString string
		err         error
	)
	expireTime := time.Duration(config.Config.TokenPolicy.AccessExpire*24) * time.Hour
	if tokenString, err = util.CreateToken(userInfo.UserID, expireTime); err != nil {
		return "", err
	}
	return tokenString, err
}

func (u *userCache) TokenInfo(token string) (string, error) {
	userId, err := util.ParseToken(token)
	if err != nil {
		return "", err
	}

	return userId, err
}

func (u *userCache) GetPrizeKey(userId string) string {
	return fmt.Sprintf("prize_%s", userId)
}

func (u *userCache) TokenVerify(token string) (status bool) {
	if _, err := util.ParseToken(token); err != nil {
		return false
	}
	return true
}

func (u *userCache) TokenRemove(token string, platformId int64) {
	userId, err := util.ParseToken(token)
	if err != nil {
		return
	}
	pId := constant.PlatformIDToName(int(platformId))
	_ = db.RedisUtil.Delete(db.TokenKey, u.getTokenIndex(userId, pId))
}

func (u *userCache) getTokenIndex(userID, sessionID string) string {
	return userID + "_" + sessionID
}

func (u *userCache) AutoDelToken(UserID string) {
	var (
		cursor uint64
		keys   []string
	)
	for {

		keys, cursor, _ = db.RedisCli.Scan(context.Background(), cursor, db.TokenKey+UserID, 20).Result()
		for _, historySession := range keys {
			_ = db.RedisCli.Del(context.Background(), historySession)
		}
		if cursor == 0 {
			break
		}
	}
}

func (u *userCache) DelUserInfoInCache(UserID string) (int64, error) {
	key := db.UserInfoKey + UserID
	return db.RedisCli.Del(context.Background(), key).Result()
}

func (u *userCache) GetUserInfoFromCache(UserID string) (*model.UserBaseInfo, error) {
	var (
		err         error
		userInfoStr string
		userInfo    *model.UserBaseInfo
	)
	key := db.UserInfoKey + UserID

	if userInfoStr, err = db.RedisCli.Get(context.Background(), key).Result(); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(userInfoStr), &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (u *userCache) DelUserInfoOnCache(UserID string) bool {
	key := db.UserInfoKey + UserID
	result, err := db.RedisCli.Del(context.Background(), key).Result()
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
	}

	return result > 0
}

func (u *userCache) GetBaseUserInfo(UserID string, fn func(UserId string) (*model.UserBaseInfo, error)) (*model.UserBaseInfo, error) {
	var (
		err   error
		setNx bool

		userInfo     *model.UserBaseInfo
		userInfoByte []byte
	)
	lockKey := db.UserInfoKey + "lock_" + UserID
	if userInfo, err = u.GetUserInfoFromCache(UserID); userInfo != nil {

		return userInfo, nil
	}
	if err != nil {
		logger.Sugar.Warn(util.GetSelfFuncName(), UserID, "userInfo from cache warn:", err)
	}

	randStr := util.RandString(12)
	expiration := time.Millisecond * 5
	if setNx, err = db.RedisCli.SetNX(context.Background(), lockKey, randStr, expiration).Result(); err != nil {
		return nil, err
	}
	defer func() {
		db.RedisCli.Del(context.Background(), lockKey)

	}()
	if !setNx {
		if userInfo, err = u.retryReadCache(UserID); err != nil {
			return nil, err
		}
		return userInfo, nil
	}

	if userInfo, err = fn(UserID); err != nil {
		return nil, err
	}
	if userInfoByte, err = json.Marshal(&userInfo); err != nil {
		return nil, err
	}
	expireTime := time.Duration(config.Config.TokenPolicy.AccessExpire*24) * time.Hour
	if _, err = db.RedisCli.SetNX(context.Background(), db.UserInfoKey+UserID, string(userInfoByte), expireTime).Result(); err != nil {
		return nil, err
	}

	if err != nil {
		logger.Sugar.Warn(util.GetSelfFuncName(), "UserInfoKey from cache warn:", err)
	}

	return userInfo, nil
}

func (u *userCache) retryReadCache(userId string) (*model.UserBaseInfo, error) {
	var (
		r        = time.Millisecond
		i        = 1
		err      = code.ErrLoadUserInfo
		userInfo *model.UserBaseInfo
	)
	for i <= RetryReadRedisTimes {
		time.Sleep(r)
		userInfo, err = u.GetUserInfoFromCache(userId)
		if userInfo != nil {
			return userInfo, nil
		}
		i++
	}

	return nil, err
}

func (u *userCache) RecordRegisterIPCount(ip string) {
	_, err := db.RedisCli.Incr(context.Background(), fmt.Sprintf("%sip_count_%s", db.UserIPRegisterKey, ip)).Result()
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
	}
	return
}

func (u *userCache) RecordRegisterDeviceIDCount(deviceID string) {
	_, err := db.RedisCli.Incr(context.Background(), fmt.Sprintf("%sdevice_count_%s", db.UserIPRegisterKey, deviceID)).Result()
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
	}
	return
}
func (u *userCache) CheckDeviceLimit(deviceID string, limit int64) bool {
	var (
		count int64
	)
	if limit == 0 {
		return false
	}
	if count, _ = db.RedisCli.Get(context.Background(), fmt.Sprintf("%sdevice_count_%s", db.UserIPRegisterKey, deviceID)).Int64(); count >= limit {
		return true
	}
	return false
}

func (u *userCache) SetIpLimitTime(ip string, t float64) bool {
	expireTime := time.Duration(math.Ceil(t*24*60)) * time.Minute
	b, _ := db.RedisCli.SetNX(context.Background(), fmt.Sprintf("%sip_limit_%s", db.UserIPRegisterKey, ip), ip, expireTime).Result()
	return b
}

func (u *userCache) CheckIpLimit(ip string, limit int64, t float64) bool {
	var (
		count      int64
		limitCount int64
	)
	if limit == 0 {

		return false
	}
	if count, _ = db.RedisCli.Get(context.Background(), fmt.Sprintf("%sip_count_%s", db.UserIPRegisterKey, ip)).Int64(); count >= limit {
		u.SetIpLimitTime(ip, t)
		u.ClearIPCount(ip)
		return true
	}
	if limitCount, _ = db.RedisCli.Exists(context.Background(), fmt.Sprintf("%sip_limit_%s", db.UserIPRegisterKey, ip)).Result(); limitCount > 0 {
		return true
	}

	return false
}

func (u *userCache) ClearIPCount(ip string) {
	db.RedisCli.Del(context.Background(), fmt.Sprintf("%sip_count_%s", db.UserIPRegisterKey, ip)).Result()
}

func (u *userCache) ClearAllIPCount() {
	ctx := context.Background()

	pattern := fmt.Sprintf("%sip_count_%s", db.UserIPRegisterKey, "*")

	var cursor uint64 = 0
	var keys []string
	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = db.RedisCli.Scan(ctx, cursor, pattern, 5000).Result()
		if err != nil {
			logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
			return
		}
		keys = append(keys, scanKeys...)
		if cursor == 0 {
			break
		}
	}

	pipe := db.RedisCli.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
		return
	}
	return
}

func (u *userCache) ClearAllIPLimit() {
	ctx := context.Background()

	pattern := fmt.Sprintf("%sip_limit_%s", db.UserIPRegisterKey, "*")

	var cursor uint64 = 0
	var keys []string
	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = db.RedisCli.Scan(ctx, cursor, pattern, 5000).Result()
		if err != nil {
			logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
			return
		}
		keys = append(keys, scanKeys...)
		if cursor == 0 {
			break
		}
	}

	pipe := db.RedisCli.Pipeline()
	for _, key := range keys {
		pipe.Del(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
		return
	}
	return
}

func (u *userCache) DelSystemConfigOnCache() bool {
	result, err := db.RedisCli.Del(context.Background(), db.SystemConfigKey).Result()
	if err != nil {
		logger.Sugar.Error(util.GetSelfFuncName(), "error:", err)
	}

	return result > 0
}

func (u *userCache) GetSystemConfigFromCache() (*cmsConfigModel.ParameterConfigResp, error) {
	var (
		err          error
		sysConfigStr string
		sysConfig    *cmsConfigModel.ParameterConfigResp
	)
	if sysConfigStr, err = db.RedisCli.Get(context.Background(), db.SystemConfigKey).Result(); sysConfigStr == "" {
		return nil, err
	}
	if err = json.Unmarshal([]byte(sysConfigStr), &sysConfig); err != nil {
		return nil, err
	}

	return sysConfig, nil
}

func (u *userCache) GetSystemConfigInfo(fn func() (*cmsConfigModel.ParameterConfigResp, error)) (*cmsConfigModel.ParameterConfigResp, error) {
	var (
		err           error
		setNx         bool
		rs            string
		sysConfig     *cmsConfigModel.ParameterConfigResp
		sysConfigByte []byte
	)
	if sysConfig, err = u.GetSystemConfigFromCache(); sysConfig != nil {
		return sysConfig, nil
	}

	randStr := util.RandString(12)
	expiration := time.Millisecond * 5
	if setNx, err = db.RedisCli.SetNX(context.Background(), db.SystemConfigLockKey, randStr, expiration).Result(); err != nil {
		return nil, err
	}
	if !setNx {
		return nil, code.ErrFailRequest
	}

	if sysConfig, err = fn(); err != nil {
		return nil, err
	}
	if sysConfigByte, err = json.Marshal(sysConfig); err != nil {
		return nil, err
	}
	expireTime := time.Duration(config.Config.TokenPolicy.AccessExpire*24) * time.Hour
	if _, err = db.RedisCli.SetNX(context.Background(), db.SystemConfigKey, string(sysConfigByte), expireTime).Result(); err != nil {
		return nil, err
	}

	if rs, err = db.RedisCli.Get(context.Background(), db.SystemConfigLockKey).Result(); rs != "" {
		if rs == randStr {
			db.RedisCli.Del(context.Background(), db.SystemConfigLockKey)
		}
	}
	if err != nil {
		logger.Sugar.Warn(util.GetSelfFuncName(), "sysConfig from cache warn:", err)
	}

	return sysConfig, nil
}

func (u *userCache) GetMemberKey(UserID string, isOwner, page interface{}) string {
	if page == "*" {
		return fmt.Sprintf("%s:%s:*", db.GetMomentsMessage, UserID)
	}
	return fmt.Sprintf("%s:%s:%d_%d", db.GetMomentsMessage, UserID, isOwner, page)
}

func (u *userCache) GetMomentsMessage(UserID string, isOwner int64, page int) (list []momentsModel.IssueInfo, err error, getCount string) {
	key := u.GetMemberKey(UserID, isOwner, page)

	data, err := db.RedisCli.Get(context.Background(), key).Result()
	if err != nil {
		return
	}
	err = util.JsonUnmarshal([]byte(data), &list)

	key = u.GetMemberKey(UserID, isOwner, UserID+"count")
	getCount, err = db.RedisCli.Get(context.Background(), key).Result()
	if err != nil {
		return
	}

	return
}

func (u *userCache) DelMomentsMessage(UserID string) {
	key := u.GetMemberKey(UserID, 0, "*")
	_ = db.RedisCli.Del(context.Background(), key)
	return
}

func (u *userCache) SetMomentsMessage(UserID string, isOwner int64, page int, list []momentsModel.IssueInfo, count int64) (err error) {
	key := u.GetMemberKey(UserID, isOwner, page)

	data, err := json.Marshal(list)
	if err != nil {
		return
	}
	err = db.RedisCli.Set(context.Background(), key, string(data), time.Minute).Err()

	key = u.GetMemberKey(UserID, isOwner, UserID+"count")
	err = db.RedisCli.Set(context.Background(), key, count, time.Minute).Err()
	return
}
