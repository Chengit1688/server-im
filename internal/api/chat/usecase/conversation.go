package usecase

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
)

var ConversationUseCase = new(conversationUseCase)

type conversationUseCase struct{}

func (c *conversationUseCase) GetConversationID(conversationType chatModel.ConversationType, userID string, recvID string) string {
	switch conversationType {
	case chatModel.ConversationTypeSingle:
		if userID == "" || recvID == "" {
			logger.Sugar.Debugw("", "func", util.GetSelfFuncName(), "info", fmt.Sprintf("user id: %s, recv id: %s", userID, recvID))
			return ""
		}

		first := util.StringToInt64(userID)
		second := util.StringToInt64(recvID)
		if first > second {
			first, second = second, first
		}
		return fmt.Sprintf("%d_%d", first, second)
	default:
		return recvID
	}
}

func (c *conversationUseCase) GetRecvID(userID string, conversationType chatModel.ConversationType, conversationID string) string {
	switch conversationType {
	case chatModel.ConversationTypeSingle:
		userIDList := strings.Split(conversationID, "_")
		for _, u := range userIDList {
			if u != userID {
				return u
			}
		}
	default:
		return conversationID
	}

	logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Errorf("parse conversation id error, conversation id: %s", conversationID))
	return ""
}

func (c *conversationUseCase) GetMaxVersion(userID string) (version int64, err error) {
	version, err = chatRepo.ConversationCache.GetVersion(userID)
	if err != nil && err != redis.Nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get version error, user id: %s, error: %v", userID, version))
		return
	}

	if err == nil {
		return
	}

	key := chatRepo.ConversationCache.GetVersionKey(userID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	if version, err = chatRepo.ConversationCache.GetVersion(userID); err == nil {
		return
	}

	if version, err = chatRepo.ConversationRepo.GetMaxVersion(userID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get max version error, key: %s, user id: %s, error: %v", key, userID, version))
		return
	}

	err = chatRepo.ConversationCache.SetVersion(userID, version)
	return
}

func (c *conversationUseCase) GetNextVersion(userID string) (version int64, err error) {
	if _, err = c.GetMaxVersion(userID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get max version error, user id: %s, error: %v", userID, version))
		return
	}

	version, err = chatRepo.ConversationCache.GetNextVersion(userID)
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get next version error, user id: %s, error: %v", userID, version))
		return
	}
	return
}

func (c *conversationUseCase) GetReadSeq(conversationType chatModel.ConversationType, conversationID string) (readSeq int64, err error) {
	readSeq, err = chatRepo.ConversationCache.GetReadSeq(conversationType, conversationID)
	if err != nil && err != redis.Nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get read seq error, conversation id: %s, error: %v", conversationID, readSeq))
		return
	}

	if err == nil {
		return
	}

	key := chatRepo.ConversationCache.GetReadSeqKey(conversationType, conversationID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	if readSeq, err = chatRepo.ConversationCache.GetReadSeq(conversationType, conversationID); err == nil {
		return
	}

	if readSeq, err = chatRepo.ConversationRepo.GetReadSeq(conversationType, conversationID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get read seq error, key: %s, conversation id: %s, error: %v", key, conversationID, readSeq))
		return
	}

	err = chatRepo.ConversationCache.SetReadSeq(conversationType, conversationID, readSeq)
	return
}

func (c *conversationUseCase) SetReadSeq(conversationType chatModel.ConversationType, conversationID string, readSeq int64) (success bool, err error) {
	key := chatRepo.ConversationCache.GetReadSeqKey(conversationType, conversationID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	var oldReadSeq int64
	if oldReadSeq, err = chatRepo.ConversationCache.GetReadSeq(conversationType, conversationID); err != nil && err != redis.Nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get read seq error, key: %s, error: %v", key, err))
		return
	}

	if readSeq > oldReadSeq {
		success = true
		err = chatRepo.ConversationCache.SetReadSeq(conversationType, conversationID, readSeq)
	}
	return
}

func (c *conversationUseCase) UpdateAckSeq(operationID string, userID string, conversationType chatModel.ConversationType, conversationID string, ackSeq int64) (err error) {

	if conversationType == chatModel.ConversationTypeSingle {
		if _, err = c.GetNextVersion(userID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get next version error, user id: %s, error: %v", userID, err))
			err = code.ErrUnknown
			return
		}
	}

	return
}

func (c *conversationUseCase) UpdateCleanSeq(operationID string, userID string, conversationType chatModel.ConversationType, conversationID string, cleanSeq int64) (err error) {

	var version int64
	if conversationType == chatModel.ConversationTypeSingle {
		if version, err = c.GetNextVersion(userID); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get next version error, user id: %s, error: %v", userID, err))
			err = code.ErrUnknown
			return
		}
	}

	err = chatRepo.ConversationRepo.UpdateCleanSeq(userID, conversationType, conversationID, cleanSeq, version)
	return
}

func (c *conversationUseCase) UpsertUsersStartSeq(operationID string, conversationType chatModel.ConversationType, conversationID string, startSeq int64, userIDList ...string) (err error) {
	var conversations []chatModel.Conversation
	for _, userID := range userIDList {
		var conversation chatModel.Conversation
		conversation.UserID = userID
		conversation.ConversationType = conversationType
		conversation.ConversationID = conversationID
		conversation.StartSeq = startSeq
		conversations = append(conversations, conversation)
	}
	err = c.upsertUsersConversations(operationID, conversations)
	return
}

func (c *conversationUseCase) upsertUsersConversations(operationID string, conversations []chatModel.Conversation) (err error) {
	for i, conversation := range conversations {
		if conversation.ConversationID == "" {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", "conversation id nil")
			err = code.ErrUnknown
			return
		}

		if conversation.ConversationType == chatModel.ConversationTypeSingle {
			if conversation.Version, err = c.GetNextVersion(conversation.UserID); err != nil {
				logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get next version error, user id: %s, error: %v", conversation.UserID, err))
				err = code.ErrUnknown
				return
			}
			conversations[i] = conversation
		}
	}

	err = chatRepo.ConversationRepo.MultiUpsert(conversations)
	return
}
