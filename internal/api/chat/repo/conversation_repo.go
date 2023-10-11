package repo

import (
	chatModel "im/internal/api/chat/model"
	"im/pkg/db"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ConversationRepo = new(conversationRepo)

type conversationRepo struct{}

func (r *conversationRepo) ListByConversationId(ConversationId string, conversationType chatModel.ConversationType, userId string) (list []chatModel.Conversation, count int64, err error) {
	var listDB *gorm.DB
	if userId != "" {
		listDB = db.DB.Model(&chatModel.Conversation{}).Where("conversation_id=? and conversation_type=? and user_id=?", ConversationId, conversationType, userId).Order("version ASC")
	} else {
		listDB = db.DB.Model(&chatModel.Conversation{}).Where("conversation_id=? and conversation_type=?", ConversationId, conversationType).Order("version ASC")
	}
	if err = listDB.Count(&count).Error; err != nil {
		return
	}
	if err = listDB.Find(&list).Error; err != nil {
		return
	}
	return
}
func (r *conversationRepo) List(userID string, version int64, offset int, limit int) (list []chatModel.Conversation, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)
	var listDB *gorm.DB
	listDB = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ? AND (version >= ? OR (conversation_type = ? AND version = 0))", userID, version, chatModel.ConversationTypeGroup).Order("version ASC")

	if needPaging {
		if err = listDB.Count(&count).Error; err != nil {
			return
		}
	}

	if needPaging {
		listDB = listDB.Offset(offset).Limit(limit)
	}

	if err = listDB.Find(&list).Error; err != nil {
		return
	}
	return
}

func (r *conversationRepo) Get(userID string, conversationType chatModel.ConversationType, conversationID string) (conversation *chatModel.Conversation, err error) {
	conversation = new(chatModel.Conversation)
	err = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ? AND conversation_type = ? AND conversation_id = ?", userID, conversationType, conversationID).
		Find(conversation).Error
	return
}

func (r *conversationRepo) GetMaxVersion(userID string) (version int64, err error) {
	data := struct {
		Version int64 `json:"version"`
	}{}

	err = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ?", userID).Order("version DESC").Limit(1).Select("version").Scan(&data).Error

	if err != nil {
		return
	}

	version = data.Version
	return
}

func (r *conversationRepo) GetReadSeq(conversationType chatModel.ConversationType, conversationID string) (readSeq int64, err error) {
	data := struct {
		AckSeq int64 `json:"ack_seq"`
	}{}

	err = db.DB.Model(&chatModel.Conversation{}).Where("conversation_type = ? AND conversation_id = ?", conversationType, conversationID).Order("ack_seq DESC").Limit(1).Select("ack_seq").Scan(&data).Error
	if err != nil {
		return
	}

	readSeq = data.AckSeq
	return
}

func (r *conversationRepo) UpdateCleanSeq(userID string, conversationType chatModel.ConversationType, conversationID string, cleanSeq int64, version int64) (err error) {
	err = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ? AND conversation_type = ? AND conversation_id = ?", userID, conversationType, conversationID).
		Updates(map[string]interface{}{"clean_seq": cleanSeq, "start_seq": cleanSeq, "version": version}).Error
	return
}

func (r *conversationRepo) UpdateVersion(userID string, conversationType chatModel.ConversationType, conversationID string, version int64) (err error) {
	err = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ? AND conversation_type = ? AND conversation_id = ?", userID, conversationType, conversationID).
		Updates(map[string]interface{}{"version": version}).Error
	return
}

func (r *conversationRepo) MultiUpsert(conversations []chatModel.Conversation) (err error) {
	for _, conversation := range conversations {
		var count int64
		if err = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ? AND conversation_type = ? AND conversation_id = ?", conversation.UserID, conversation.ConversationType, conversation.ConversationID).
			Count(&count).Error; err != nil {
			return err
		}

		if count < 1 {

			if err = db.DB.Model(&chatModel.Conversation{}).Create(&conversation).Error; err != nil {
				return err
			}
		} else {

			if err = db.DB.Model(&chatModel.Conversation{}).Where("user_id = ? AND conversation_type = ? AND conversation_id = ?", conversation.UserID, conversation.ConversationType, conversation.ConversationID).
				Updates(map[string]interface{}{"deleted_at": 0, "start_seq": conversation.StartSeq, "version": conversation.Version}).Error; err != nil {
				return err
			}
		}
	}
	return
}

func (r *conversationRepo) Delete(conversationType chatModel.ConversationType, conversationID string, userIDList ...string) (err error) {
	var query string
	var args []interface{}

	if len(userIDList) != 0 {
		query += " AND user_id IN (?)"
		args = append(args, userIDList)
	}

	query += " AND conversation_type = ? AND conversation_id = ?"
	args = append(args, conversationType, conversationID)

	query = strings.TrimPrefix(strings.TrimSpace(query), "AND")

	err = db.DB.Model(&chatModel.Conversation{}).Where(query, args...).Updates(map[string]interface{}{"deleted_at": time.Now().Unix()}).Error
	return
}
