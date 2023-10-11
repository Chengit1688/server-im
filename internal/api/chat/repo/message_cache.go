package repo

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	chatModel "im/internal/api/chat/model"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
)

var MessageCache = new(messageCache)

type messageCache struct{}

func (c *messageCache) GetStatusKey(status chatModel.MessageStatusType, msgID string) string {
	var name string
	switch status {
	case chatModel.MessageStatusTypeRead:
		name = "read"

	case chatModel.MessageStatusTypeRevoke:
		name = "revoke"

	case chatModel.MessageStatusTypeDelete:
		name = "delete"
	}
	return fmt.Sprintf("%s:%s:%s", db.MessageStatus, name, msgID)
}

func (c *messageCache) GetStatusExist(status chatModel.MessageStatusType, msgID string) bool {
	key := c.GetStatusKey(status, msgID)

	_, err := db.RedisCli.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Errorf("redis setnx error, key: %s, error: %v", key, err))
		return true
	}

	if err == nil {
		return true
	}
	return false
}

func (c *messageCache) SetStatus(status chatModel.MessageStatusType, msgID string) {
	key := c.GetStatusKey(status, msgID)

	if _, err := db.RedisCli.SetNX(context.Background(), key, 1, util.RandDuration(util.OneMonth*6)).Result(); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Errorf("redis setnx error, key: %s, error: %v", key, err))
	}
	return
}

func (c *messageCache) GetSeqKey(conversationType chatModel.ConversationType, conversationID string) string {
	var name string
	switch conversationType {
	case chatModel.ConversationTypeSingle:
		name = "single"
	case chatModel.ConversationTypeGroup:
		name = "group"
	case chatModel.ConversationTypeChannel:
		name = "channel"
	}
	return fmt.Sprintf("%s:%s:%s", db.MessageSeq, name, conversationID)
}
