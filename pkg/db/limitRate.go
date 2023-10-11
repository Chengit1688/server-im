package db

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

func RateLimit(key,value string, limit int, expired  time.Duration) (bool, error){
	length, err := RedisCli.LLen(context.Background(),key).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	leng:=int(length)
	if leng >= limit {
		return true, nil
	}
	if leng ==0{
		RedisCli.RPush(context.Background(),key,value)
		RedisCli.Expire(context.Background(),key,expired)
	}else{
		RedisCli.RPushX(context.Background(),key,value)
	}
	return false, err
}
