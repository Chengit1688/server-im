package repo

import (
	"context"
	"fmt"
	chatModel "im/internal/api/chat/model"
	"im/pkg/db"
)

var ConversationCache = new(conversationCache)

type conversationCache struct{}

func (c *conversationCache) GetVersionKey(userID string) string {
	return fmt.Sprintf("%s:%s", db.ConversationVersion, userID)
}

func (c *conversationCache) GetVersion(userID string) (version int64, err error) {
	key := c.GetVersionKey(userID)
	version, err = db.RedisCli.Get(context.Background(), key).Int64()
	return
}

func (c *conversationCache) SetVersion(userID string, value interface{}) (err error) {
	key := c.GetVersionKey(userID)
	err = db.RedisCli.SetNX(context.Background(), key, value, 0).Err()
	return
}

func (c *conversationCache) GetNextVersion(userID string) (nextVersion int64, err error) {
	key := c.GetVersionKey(userID)
	nextVersion, err = db.RedisCli.Incr(context.Background(), key).Result()
	return
}

func (c *conversationCache) GetReadSeqKey(conversationType chatModel.ConversationType, conversationID string) string {
	var name string
	switch conversationType {
	case chatModel.ConversationTypeSingle:
		name = "single"
	case chatModel.ConversationTypeGroup:
		name = "group"
	}
	return fmt.Sprintf("%s:%s:%s", db.ConversationReadSeq, name, conversationID)
}

func (c *conversationCache) GetReadSeq(conversationType chatModel.ConversationType, conversationID string) (readSeq int64, err error) {
	key := c.GetReadSeqKey(conversationType, conversationID)
	readSeq, err = db.RedisCli.Get(context.Background(), key).Int64()
	return
}

func (c *conversationCache) SetReadSeq(conversationType chatModel.ConversationType, conversationID string, value interface{}) (err error) {
	key := c.GetReadSeqKey(conversationType, conversationID)
	err = db.RedisCli.Set(context.Background(), key, value, 0).Err()
	return
}
