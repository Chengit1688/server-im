package repo

import (
	"fmt"
	"gorm.io/gorm"
	chatModel "im/internal/api/chat/model"
	"im/pkg/db"
	"im/pkg/util"
	"strings"
)

var MessageRepo = new(messageRepo)

type messageRepo struct{}

func (r *messageRepo) PullNew(conversationType chatModel.ConversationType, conversationID string, startSeq int64, seq int64, pageSize int) (list []chatModel.MessageInfo, err error) {
	if pageSize == 0 {
		pageSize = 20
	}

	if pageSize > 999 {
		pageSize = 999
	}
	if seq < startSeq {
		seq = startSeq
	}
	var messageList []chatModel.Message
	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq > ? ", conversationType, conversationID, startSeq).
		Order("seq DESC").Limit(pageSize).Find(&messageList).Error
	util.CopyStructFields(&list, &messageList)
	return
}

func (r *messageRepo) PullNewV2(conversationType chatModel.ConversationType, conversationID string, startSeq int64, seq int64, pageSize int) (list []chatModel.MessageInfo, err error) {
	if pageSize == 0 {
		pageSize = 20
	}

	if pageSize > 999 {
		pageSize = 999
	}
	if seq < startSeq {
		seq = startSeq
	}
	var messageList []chatModel.Message
	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq <= ? ", conversationType, conversationID, seq).
		Order("seq DESC").Limit(pageSize).Find(&messageList).Error
	util.CopyStructFields(&list, &messageList)
	return
}

func (r *messageRepo) PullNewList(conversationType chatModel.ConversationType, conversationID string, seqList []int64, pageSize int) (list []chatModel.MessageInfo, err error) {
	if pageSize == 0 {
		pageSize = 20
	}

	if pageSize > 999 {
		pageSize = 999
	}

	var messageList []chatModel.Message
	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq IN ? ", conversationType, conversationID, seqList).
		Order("seq asc").Limit(pageSize).Find(&messageList).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
	util.CopyStructFields(&list, &messageList)
	return
}

func (r *messageRepo) Pull(conversationType chatModel.ConversationType, conversationID string, startSeq int64, seq int64, pageSize int) (list []chatModel.MessageInfo, err error) {
	if pageSize == 0 {
		pageSize = 20
	}

	if pageSize > 999 {
		pageSize = 999
	}

	var messageList []chatModel.Message

	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq > ? AND  seq <= ?", conversationType, conversationID, startSeq, seq).
		Order("seq DESC").Limit(pageSize).Find(&messageList).Error
	util.CopyStructFields(&list, &messageList)
	return
}

func (r *messageRepo) ListSingle(conversationID string, sendID string, recvID string, messageType chatModel.MessageType, content string, startTime int64, endTime int64, offset int, limit int) (list []chatModel.MessageInfo, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)
	var query string
	var args []interface{}

	query += " AND conversation_type = ?"
	args = append(args, chatModel.ConversationTypeSingle)

	if sendID != "" && recvID != "" {
		query += " AND conversation_id = ? AND send_id = ?"
		args = append(args, conversationID, sendID)
	} else if sendID == "" && recvID != "" {
		query += " AND conversation_id LIKE ? AND send_id != ?"
		args = append(args, fmt.Sprintf("%%%s%%", recvID), recvID)
	} else if sendID != "" && recvID == "" {
		query += " AND send_id = ?"
		args = append(args, sendID)
	}

	if messageType != chatModel.MessageNone {
		query += " AND type = ?"
		args = append(args, messageType)
	} else {
		query += " AND type < ?"
		args = append(args, chatModel.MessageOperation)
	}

	if startTime > 0 {
		query += " AND send_time >= ?"
		args = append(args, startTime)
	}

	if endTime > 0 {
		query += " AND send_time <= ?"
		args = append(args, endTime)
	}

	query += " AND status IN (?, ?)"
	args = append(args, chatModel.MessageStatusTypeNotRead, chatModel.MessageStatusTypeRead)

	if content != "" {
		query += " AND CASE WHEN JSON_VALID(content) THEN JSON_EXTRACT(content, '$.text') LIKE ? END"
		args = append(args, fmt.Sprintf("%%%s%%", content))
	}

	query = strings.TrimPrefix(strings.TrimSpace(query), "AND")

	var listDB *gorm.DB
	listDB = db.DB.Model(&chatModel.Message{}).Order("send_time DESC")

	if needPaging {
		listDB = listDB.Offset(offset).Limit(limit)
	}

	var messageList []chatModel.Message
	if err = listDB.Where(query, args...).Find(&messageList).Error; err != nil {
		return
	}

	if needPaging {
		if err = db.DB.Model(&chatModel.Message{}).Where(query, args...).Count(&count).Error; err != nil {
			return
		}
	}
	util.CopyStructFields(&list, &messageList)
	return
}

func (r *messageRepo) List(conversationID string, sendID string, messageType chatModel.MessageType, content string, startTime int64, endTime int64, offset int, limit int) (list []chatModel.MessageInfo, count int64, err error) {
	needPaging := !(offset == limit && limit == 0)
	var query string
	var args []interface{}

	query += " AND conversation_type = ?"
	args = append(args, chatModel.ConversationTypeGroup)

	if conversationID != "" {
		query += " AND conversation_id = ?"
		args = append(args, conversationID)
	}

	if sendID != "" {
		query += " AND send_id = ?"
		args = append(args, sendID)
	}

	if messageType != chatModel.MessageNone {
		query += " AND type = ?"
		args = append(args, messageType)
	} else {
		query += " AND type < ?"
		args = append(args, chatModel.MessageOperation)
	}

	if startTime > 0 {
		query += " AND send_time >= ?"
		args = append(args, startTime)
	}

	if endTime > 0 {
		query += " AND send_time <= ?"
		args = append(args, endTime)
	}

	query += " AND status IN (?, ?)"
	args = append(args, chatModel.MessageStatusTypeNotRead, chatModel.MessageStatusTypeRead)

	if content != "" {
		query += " AND CASE WHEN JSON_VALID(content) THEN JSON_EXTRACT(content, '$.text') LIKE ? END"
		args = append(args, fmt.Sprintf("%%%s%%", content))
	}

	query = strings.TrimPrefix(strings.TrimSpace(query), "AND")

	var listDB *gorm.DB
	listDB = db.DB.Model(&chatModel.Message{}).Order("send_time DESC")

	if needPaging {
		listDB = listDB.Offset(offset).Limit(limit)
	}

	var messageList []chatModel.Message
	if err = listDB.Where(query, args...).Find(&messageList).Error; err != nil {
		return
	}

	if needPaging {
		if err = db.DB.Model(&chatModel.Message{}).Where(query, args...).Count(&count).Error; err != nil {
			return
		}
	}
	util.CopyStructFields(&list, &messageList)
	return
}

func (r *messageRepo) Get(msgID string) (message *chatModel.Message, err error) {
	message = new(chatModel.Message)
	err = db.DB.Model(&chatModel.Message{}).Where("msg_id = ?", msgID).Find(message).Error
	return
}

func (r *messageRepo) Add(message *chatModel.Message) (err error) {
	err = db.DB.Model(&chatModel.Message{}).Create(message).Error
	return
}

func (r *messageRepo) UnreadCount(conversationType chatModel.ConversationType, conversationID string, seq int64, userID string) (unreadCount int64, err error) {
	var typeList []chatModel.MessageType
	for i := chatModel.MessageText; i <= chatModel.MessageQuote; i++ {
		typeList = append(typeList, i)
	}

	typeList = append(typeList, chatModel.MessageGroupNotifyChangeNotify)

	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq > ? AND type IN (?) AND send_id != ? AND status IN (?, ?)",
		conversationType, conversationID, seq, typeList, userID, chatModel.MessageStatusTypeNotRead, chatModel.MessageStatusTypeRead).Count(&unreadCount).Error
	return
}

func (r *messageRepo) LatestMessage(conversationType chatModel.ConversationType, conversationID string, seq int64) (messageInfo *chatModel.MessageInfo, err error) {
	message := new(chatModel.Message)
	messageInfo = new(chatModel.MessageInfo)
	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq > ?", conversationType, conversationID, seq).
		Order("seq DESC").Limit(1).Find(message).Error
	util.CopyStructFields(messageInfo, message)
	return
}

func (r *messageRepo) UpdateStatus(status chatModel.MessageStatusType, msgIDList []string) (err error) {
	err = db.DB.Model(&chatModel.Message{}).Where("msg_id IN (?)", msgIDList).Updates(map[string]interface{}{"status": status}).Error
	return
}

func (r *messageRepo) UpdateStatusRead(conversationType chatModel.ConversationType, conversationID string, seq int64) (err error) {
	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND seq <= ? AND status IN (?)", conversationType, conversationID, seq, chatModel.MessageStatusTypeNotRead).
		Updates(map[string]interface{}{"status": chatModel.MessageStatusTypeRead}).Error
	return
}

func (r *messageRepo) UpdateSenderMessageStatus(conversationType chatModel.ConversationType, conversationID string, status chatModel.MessageStatusType, sendIDList ...string) (err error) {
	err = db.DB.Model(&chatModel.Message{}).Where("conversation_type = ? AND conversation_id = ? AND send_id IN (?)", conversationType, conversationID, sendIDList).Updates(map[string]interface{}{"status": status}).Error
	return
}

func (r *messageRepo) Clear(clearTime int64) (count int64, err error) {
	if err = db.DB.Model(&chatModel.Message{}).Where("send_time < ?", clearTime).Count(&count).Error; err != nil {
		return
	}

	if count == 0 {
		return
	}

	err = db.DB.Where("send_time < ?", clearTime).Limit(100).Delete(&chatModel.Message{}).Error
	return
}
