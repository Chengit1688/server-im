package repo

import (
	"context"
	"fmt"
	chatModel "im/internal/api/chat/model"
	"im/pkg/db"
	"im/pkg/util"
	"time"
)

var RTCCache = new(rtcCache)

type rtcCache struct{}

func (c *rtcCache) GetRTCKey(userID string) string {
	return fmt.Sprintf("%s:%s", db.RTC, userID)
}

func (c *rtcCache) GetRTC(userID string) (rtcInfo *chatModel.RTCInfo, err error) {
	key := c.GetRTCKey(userID)

	var dataStr string
	if dataStr, err = db.RedisCli.Get(context.Background(), key).Result(); err != nil {
		return
	}

	rtcInfo = new(chatModel.RTCInfo)
	if err = util.JsonUnmarshal([]byte(dataStr), rtcInfo); err != nil {
		return
	}

	switch rtcInfo.RTCStatus {
	case chatModel.RTCStatusTypeRequest:
		rtcInfo.RTCRequestLimitTime = chatModel.RTCRequestExpireTime - (time.Now().Unix() - rtcInfo.RTCStartTime)
		if rtcInfo.RTCRequestLimitTime < 0 {
			rtcInfo.RTCRequestLimitTime = 0
		}

	case chatModel.RTCStatusTypeAgree:
		rtcInfo.RTCRetainTime = time.Now().Unix() - rtcInfo.RTCStartTime
	}
	return
}

func (c *rtcCache) SetRTC(userID string, rtcInfo interface{}, ttl int) (err error) {
	key := c.GetRTCKey(userID)

	var data []byte
	if data, err = util.JsonMarshal(&rtcInfo); err != nil {
		return
	}

	err = db.RedisCli.Set(context.Background(), key, string(data), time.Duration(ttl)*time.Second).Err()
	return
}

func (c *rtcCache) DeleteRTC(keys ...string) {
	if len(keys) == 0 {
		return
	}

	_ = db.RedisCli.Del(context.Background(), keys...).Err()
	return
}

func (c *rtcCache) GetRTCGroupKey(groupID string) string {
	return fmt.Sprintf("%s:group:%s", db.RTC, groupID)
}
