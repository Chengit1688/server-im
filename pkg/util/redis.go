package util

import (
	"context"
	"errors"
	"fmt"
	"im/pkg/logger"
	"time"

	"github.com/go-redis/redis/v9"
	jsoniter "github.com/json-iterator/go"
)

type RedisUtil struct {
	client *redis.Client
}

func NewRedisUtil(client *redis.Client) *RedisUtil {
	return &RedisUtil{client: client}
}

// Set 将指定值设置到redis中，使用json的序列化方式
func (u *RedisUtil) Set(key string, value interface{}, duration time.Duration) error {
	bytes, err := jsoniter.Marshal(value)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	err = u.client.Set(context.Background(), key, bytes, duration).Err()
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// Get 从redis中读取指定值，使用json的反序列化方式
func (u *RedisUtil) Get(key string, value interface{}) error {
	bytes, err := u.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(bytes, value)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// SetByExpiration
func (u *RedisUtil) SetByExpiration(kind string, name, resource string, expiration time.Duration) error {
	resourceKey := fmt.Sprintf("%s%s", kind, name)
	result := u.client.Set(context.Background(), resourceKey, resource, expiration*time.Second)
	if result.Err() != nil {
		return fmt.Errorf("%s(%s)写入redis错误:%s", kind, name, result.Err().Error())
	}

	if result.Val() != "OK" {
		return fmt.Errorf("%s(%s)写入redis失败", kind, name)
	}

	return nil
}

// Delete
func (u *RedisUtil) Delete(kind string, name string) error {
	resourceKey := GetResKey(kind, name)
	result := u.client.Del(context.Background(), resourceKey)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func GetResKey(kind string, name string) string {
	//schedule:groups:group111
	pre := GetResKeyPrefix(kind)
	return fmt.Sprintf("%s%s", pre, name)
}

func GetResKeyPrefix(kind string) string {
	return fmt.Sprintf("im2:%s:", kind)
}

func (u *RedisUtil) GetToken(kind string, name string) (string, error) {
	//resourceKey := GetResKey(kind, name)
	resourceKey := fmt.Sprintf("%s%s", kind, name)
	result := u.client.Get(context.Background(), resourceKey)
	if result.Err() != nil {
		if result.Err().Error() == "redis: nil" {
			return "", errors.New("not found")
		} else {
			return "", fmt.Errorf("%s(%s)读取redis错误:%s", kind, name, result.Err().Error())
		}
	}

	return result.Val(), nil
}

func (u *RedisUtil) HSet(key, field string, value interface{}) error {
	return u.client.HSet(context.Background(), key, field, value).Err()
}

const (
	unlockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`
)

var (
	errLockFailed = errors.New("redis lock failed")
	errTimeout    = errors.New("redis lock timeout")
)

type RedisLock struct {
	client          *redis.Client
	script          *redis.Script
	watchDog        chan struct{}
	key             string
	randomValue     string        // 随机值
	ttl             time.Duration // 过期时间
	tryLockMaxTime  time.Duration // 获取锁最大重试时间
	tryLockInterval time.Duration // 重新获取锁时间间隔
}

func NewLock(client *redis.Client, key string) *RedisLock {
	return &RedisLock{
		client:          client,
		script:          redis.NewScript(unlockScript),
		watchDog:        make(chan struct{}),
		key:             fmt.Sprintf("redis_lock:%s", key),
		randomValue:     fmt.Sprintf("%d_%d", time.Now().Nanosecond(), r.Intn(10000)),
		ttl:             time.Second * 5,
		tryLockMaxTime:  time.Second * 10,
		tryLockInterval: time.Millisecond * 100,
	}
}

func (l *RedisLock) IsLock() bool {
	res, err := l.client.Get(context.Background(), l.key).Result()
	if err == redis.Nil {
		return false
	}

	if res != "" {
		return true
	}
	return false
}

func (l *RedisLock) Lock() (err error) {
	if err = l.tryLock(); err == nil {
		return
	}

	if !errors.Is(err, errLockFailed) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.tryLockMaxTime)
	defer cancel()

	ticker := time.NewTicker(l.tryLockInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err = l.tryLock(); err == nil {
				return
			}

			if !errors.Is(err, errLockFailed) {
				return
			}

		case <-ctx.Done():
			err = errTimeout
			return
		}
	}
}

func (l *RedisLock) Unlock() {
	_ = l.script.Run(context.Background(), l.client, []string{l.key}, l.randomValue).Err()
	close(l.watchDog)
}

func (l *RedisLock) tryLock() (err error) {
	var ok bool
	if ok, err = l.client.SetNX(context.Background(), l.key, l.randomValue, l.ttl).Result(); err != nil {
		return
	}

	if !ok {
		err = errLockFailed
		return
	}

	go l.startWatchDog()
	return
}

// 看门狗，延长锁持有时间
func (l *RedisLock) startWatchDog() {
	ticker := time.NewTicker(l.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-l.watchDog:
			return
		case <-ticker.C:
			// 延长锁持有时间，失败或不存在返回
			ok, err := l.client.Expire(context.Background(), l.key, l.ttl).Result()
			if err != nil || !ok {
				return
			}
		}
	}
}
