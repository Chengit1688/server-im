package usecase

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	"im/internal/api/group/model"
	"im/internal/api/group/repo"
	userModel "im/internal/api/user/model"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/minio"
	"im/pkg/mqtt"
	"im/pkg/util"
	"strings"
	"time"
)

var MessageUseCase = new(messageUseCase)

type messageUseCase struct{}

func (c *messageUseCase) SendMessageToUsers(operationID string, clientMsgID string, conversationType chatModel.ConversationType, conversationID string, sendID string, messageType chatModel.MessageType, content string, userIDList ...string) (message *chatModel.MessageInfo, err error) {
	if message, err = c.saveMessage(operationID, clientMsgID, conversationType, conversationID, sendID, messageType, content); err != nil {
		return
	}

	switch conversationType {
	case chatModel.ConversationTypeSingle:

		if len(userIDList) != 2 {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user id list size error, size: %d", len(userIDList)))
			return
		}

		userID1, userID2 := userIDList[0], userIDList[1]
		switch sendID {
		case userID1:
			message.RecvID = userID2

		case userID2:
			message.RecvID = userID1

		default:
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send message to users error, send id not match, send id: %s, user id list: %v", sendID, userIDList))
		}

		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessagePush, message, userIDList...)

	case chatModel.ConversationTypeGroup:

		message.RecvID = conversationID
		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessagePush, message, userIDList...)
	}
	return
}

func (c *messageUseCase) SendMessageToGroup(operationID string, clientMsgID string, conversationType chatModel.ConversationType, conversationID string, sendID string, messageType chatModel.MessageType, content string) (message *chatModel.MessageInfo, err error) {
	if message, err = c.saveMessage(operationID, clientMsgID, conversationType, conversationID, sendID, messageType, content); err != nil {
		return
	}

	message.RecvID = conversationID
	_ = mqtt.SendMessageToGroups(operationID, common.ChatMessagePush, message, conversationID)
	return
}

func (c *messageUseCase) GetMember(groupID string, memberID string) (member *model.GroupMember, err error) {
	member, err = repo.GroupMemberCache.GetMember(groupID, memberID)
	if err != nil && err != redis.Nil {
		return
	}

	if err == nil {
		return
	}

	key := repo.GroupMemberCache.GetMemberKey(groupID, memberID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	if member, err = repo.GroupMemberCache.GetMember(groupID, memberID); err == nil {
		return
	}

	if member, err = repo.GroupMemberRepo.GetMember(groupID, memberID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get member error, error: %v", err))
		return
	}

	if member.Id == 0 {
		err = errors.New("record not found")
		return
	}

	err = repo.GroupMemberCache.SetMember(groupID, member)
	return
}

func (c *messageUseCase) saveMessage(operationID string, clientMsgID string, conversationType chatModel.ConversationType, conversationID string, sendID string, messageType chatModel.MessageType, content string) (message *chatModel.MessageInfo, err error) {
	var (
		user           *userModel.UserBaseInfo
		decryptContent string
		nextSeq        int64
	)

	if conversationID == "" {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation id nil"))
		err = code.ErrUnknown
		return
	}

	if user, err = userUseCase.UserUseCase.GetBaseInfo(sendID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, user id: %s, error: %v", sendID, err))
		err = code.ErrDB
		return
	}

	if decryptContent, err = util.Decrypt(content, common.ContentKey); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("decrypt error, error: %v", err))
		err = code.ErrUnknown
		return
	}

	realNickname := user.NickName

	var sendRole model.RoleType
	groupMember, err := repo.GroupMemberRepo.GetMember(conversationID, sendID)
	if err == nil && groupMember != nil {
		realNickname = groupMember.GroupNickName
		if realNickname == "" {
			realNickname = user.NickName
		}
		sendRole = groupMember.Role
	}
	message = new(chatModel.MessageInfo)
	message.ClientMsgID = clientMsgID
	message.ConversationType = conversationType
	message.ConversationID = conversationID
	message.SendID = sendID
	message.SendNickname = realNickname
	message.SendFaceUrl = user.FaceURL
	message.SendTime = util.UnixMilliTime(time.Now())
	message.Type = messageType
	message.Content = decryptContent
	message.Seq = nextSeq
	message.MsgID = c.GetMsgID(conversationType, conversationID, message.SendTime)
	message.Role = string(sendRole)

	if message.ClientMsgID == "" {
		message.ClientMsgID = message.MsgID
	}
	logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf(" group push message %+v", message))

	if err = chatRepo.MessageRepo.Add(&message.Message); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("add error, error: %v", err))
		err = code.ErrDB
		return
	}

	message.Content = content
	return
}

func (c *messageUseCase) SendSystemMessageToUsers(operationID string, conversationType chatModel.ConversationType, conversationID string, messageType chatModel.MessageType, content *chatModel.MessageContent, userIDList ...string) (message *chatModel.MessageInfo, err error) {
	if message, err = c.saveSystemMessage(operationID, conversationType, conversationID, messageType, content); err != nil {
		return
	}

	switch conversationType {
	case chatModel.ConversationTypeSingle:

		if len(userIDList) != 2 {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user id list size error, size: %d", len(userIDList)))
			return
		}

		userID1, userID2 := userIDList[0], userIDList[1]
		userID1Info, _ := userUseCase.UserUseCase.GetBaseInfo(userID1)
		message.RecvID = userID1
		message.Message.SendID = userID1
		message.SendFaceUrl = userID1Info.FaceURL
		message.SendNickname = userID1Info.NickName
		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessagePush, message, userID2)

		userID2Info, _ := userUseCase.UserUseCase.GetBaseInfo(userID2)
		message.RecvID = userID2
		message.Message.SendID = userID2
		message.SendFaceUrl = userID2Info.FaceURL
		message.SendNickname = userID2Info.NickName
		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessagePush, message, userID1)

	case chatModel.ConversationTypeGroup:

		message.RecvID = conversationID
		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessagePush, message, userIDList...)
	}
	return
}

func (c *messageUseCase) SendSystemMessageToGroup(operationID string, conversationID string, messageType chatModel.MessageType, content *chatModel.MessageContent) (message *chatModel.MessageInfo, err error) {
	if message, err = c.saveSystemMessage(operationID, chatModel.ConversationTypeGroup, conversationID, messageType, content); err != nil {
		return
	}

	message.RecvID = conversationID
	_ = mqtt.SendMessageToGroups(operationID, common.ChatMessagePush, message, conversationID)
	return
}

func (c *messageUseCase) SendSystemMessageToGroupAndGroupMembers(operationID string, conversationID string, messageType chatModel.MessageType, content *chatModel.MessageContent, groupMemberIDList ...string) (message *chatModel.MessageInfo, err error) {
	if message, err = c.saveSystemMessage(operationID, chatModel.ConversationTypeGroup, conversationID, messageType, content); err != nil {
		return
	}

	message.RecvID = conversationID

	_ = mqtt.SendMessageToUsers(operationID, common.ChatMessagePush, message, groupMemberIDList...)
	_ = mqtt.SendMessageToGroups(operationID, common.ChatMessagePush, message, conversationID)
	return
}

func (c *messageUseCase) FillSystemMessageUserInfo(operationID string, message *chatModel.MessageInfo) {
	if message.Type < chatModel.MessageOperation {
		return
	}

	var content chatModel.MessageContent
	if err := util.JsonUnmarshal([]byte(message.Content), &content); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("json unmarshal error, error: %v", err))
		return
	}

	if content.OperatorID != "" {
		operator, err := userUseCase.UserUseCase.GetBaseInfo(content.OperatorID)
		if err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, operator id: %s, error: %v", content.OperatorID, err))
			return
		}
		content.OperatorNickname = operator.NickName
		content.OperatorFaceUrl = operator.FaceURL
	}

	for i, beOperatorInfo := range content.BeOperatorList {
		beOperator, err := userUseCase.UserUseCase.GetBaseInfo(beOperatorInfo.BeOperatorID)
		if err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, be operator id: %s, error: %v", beOperatorInfo.BeOperatorID, err))
			return
		}
		content.BeOperatorList[i].BeOperatorNickname = beOperator.NickName
	}

	data, err := util.JsonMarshal(&content)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("json marshal error, error: %v", err))
		return
	}

	message.Content = string(data)
}

func (c *messageUseCase) saveSystemMessage(operationID string, conversationType chatModel.ConversationType, conversationID string, messageType chatModel.MessageType, content *chatModel.MessageContent) (message *chatModel.MessageInfo, err error) {
	if conversationID == "" {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation id nil"))
		err = code.ErrUnknown
		return
	}

	var nextSeq int64
	var data []byte
	if data, err = util.JsonMarshal(content); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("json marshal error, error: %v", err))
		err = code.ErrUnknown
		return
	}

	message = new(chatModel.MessageInfo)
	message.ConversationType = conversationType
	message.ConversationID = conversationID
	message.SendTime = util.UnixMilliTime(time.Now())
	message.Type = messageType
	message.Content = string(data)
	message.Status = chatModel.MessageStatusTypeRead
	message.Seq = nextSeq
	message.MsgID = c.GetMsgID(conversationType, conversationID, message.SendTime)
	message.ClientMsgID = message.MsgID

	if err = chatRepo.MessageRepo.Add(&message.Message); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("add error, error: %v", err))
		err = code.ErrDB
		return
	}

	c.FillSystemMessageUserInfo(operationID, message)

	if message.Content, err = util.Encrypt([]byte(message.Content), common.ContentKey); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("util encrypt error, error: %v", err))
		err = code.ErrUnknown
		return
	}
	return
}

func (c *messageUseCase) UpdateStatus(operationID string, userID string, conversationType chatModel.ConversationType, conversationID string, messageType chatModel.MessageType, status chatModel.MessageStatusType, msgIDList []string) (message *chatModel.MessageInfo, err error) {

	var realMsgIDList []string
	for _, msgID := range msgIDList {
		if !chatRepo.MessageCache.GetStatusExist(status, msgID) {
			realMsgIDList = append(realMsgIDList, msgID)
		}
	}

	if len(realMsgIDList) == 0 {
		return
	}

	if err = chatRepo.MessageRepo.UpdateStatus(status, realMsgIDList); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("update status error, error: %v", err))
		err = code.ErrDB
		return
	}

	var content chatModel.MessageContent
	content.OperatorID = userID
	content.MsgIDList = realMsgIDList

	switch conversationType {
	case chatModel.ConversationTypeSingle:
		userIDList := strings.Split(conversationID, "_")
		if message, err = c.SendSystemMessageToUsers(operationID, conversationType, conversationID, messageType, &content, userIDList...); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("send system message to users error, error: %v", err))
			return
		}

	case chatModel.ConversationTypeGroup:
		if message, err = c.SendSystemMessageToGroup(operationID, conversationID, messageType, &content); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("send system message to group error, error: %v", err))
			return
		}
	}

	for _, msgID := range realMsgIDList {
		chatRepo.MessageCache.SetStatus(status, msgID)
	}
	return
}

func (c *messageUseCase) GetMsgID(conversationType chatModel.ConversationType, conversationID string, sendTime int64) (name string) {
	randNum := util.RandInt(0, sendTime)

	switch conversationType {
	case chatModel.ConversationTypeSingle:
		name = fmt.Sprintf("single_%s_%d_%d", conversationID, sendTime, randNum)
	case chatModel.ConversationTypeGroup:
		name = fmt.Sprintf("group_%s_%d_%d", conversationID, sendTime, randNum)
	}
	return util.Md5(name)
}

func (c *messageUseCase) ComposeMessageContent(operationID string, contentStr string) (content string) {
	var data chatModel.MessageContentClient
	if err := util.JsonUnmarshal([]byte(contentStr), &data); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("json unmarshal error, error: %v", err))
		return
	}

	if data.Text != "" {
		content += fmt.Sprintf(" %s", data.Text)
	}

	if data.ImageInfo.ImageURL != "" {
		content += fmt.Sprintf(" %s", minio.GetRealURL(data.ImageInfo.ImageURL))
	}

	if data.AudioInfo.FileURL != "" {
		content += fmt.Sprintf(" %s", minio.GetRealURL(data.AudioInfo.FileURL))
	}

	if data.VideoInfo.FileURL != "" {
		content += fmt.Sprintf(" %s", minio.GetRealURL(data.VideoInfo.FileURL))
	}

	if data.FileInfo.FileURL != "" {
		content += fmt.Sprintf(" %s", minio.GetRealURL(data.FileInfo.FileURL))
	}

	if data.CardInfo.Nickname != "" {
		content += fmt.Sprintf(" %s", data.CardInfo.Nickname)
	}

	content = strings.TrimPrefix(content, " ")
	return
}

func (c *messageUseCase) ClearClientNew(operationID string, userID string, conversationType chatModel.ConversationType, conversationID string) (err error) {

	var maxSeq int64
	if err = ConversationUseCase.UpdateCleanSeq(operationID, userID, conversationType, conversationID, maxSeq); err != nil {
		return
	}

	return
}

func (c *messageUseCase) ClearClient(operationID string, userID string) (err error) {
	var conversations []chatModel.Conversation
	if conversations, _, err = chatRepo.ConversationRepo.List(userID, 0, 0, 0); err != nil {
		return
	}

	for i, conversation := range conversations {
		var maxSeq int64
		if conversation.DeletedAt != 0 || conversation.CleanSeq == maxSeq {
			continue
		}

		if err = ConversationUseCase.UpdateCleanSeq(operationID, userID, conversation.ConversationType, conversation.ConversationID, maxSeq); err != nil {
			return
		}

		if i != 0 && i%50 == 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
	return
}

func (c *messageUseCase) ClearConversation(operationID string, userID string, conversationType chatModel.ConversationType, conversationID string, maxSeq int64) (err error) {
	var conversation *chatModel.Conversation
	if conversation, err = chatRepo.ConversationRepo.Get(userID, conversationType, conversationID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get error, conversation type: %d, conversation id: %s, error: %v", conversationType, conversationID, err))
		err = code.ErrDB
		return
	}

	if conversation.ID == 0 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation not found, conversation type: %d, conversation id: %s", conversationType, conversationID))
		err = code.ErrUnknown
		return
	}

	if maxSeq < conversation.StartSeq {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("max seq < start seq, max: %d, start seq: %d", maxSeq, conversation.StartSeq))
		err = code.ErrBadRequest
		return
	}

	if err = ConversationUseCase.UpdateCleanSeq(operationID, userID, conversationType, conversationID, maxSeq); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update clean seq error, conversation type: %d, conversation id: %s, error: %v", conversationType, conversationID, err))
		err = code.ErrDB
		return
	}
	return
}

func (c *messageUseCase) ClearConversationNew(operationID string, userID string, conversationType chatModel.ConversationType, conversationID string, maxSeq int64) (err error) {
	var conversation *chatModel.Conversation
	if conversation, err = chatRepo.ConversationRepo.Get(userID, conversationType, conversationID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get error, conversation type: %d, conversation id: %s, error: %v", conversationType, conversationID, err))
		err = code.ErrDB
		return
	}

	if conversation.ID == 0 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation not found, conversation type: %d, conversation id: %s", conversationType, conversationID))
		err = code.ErrUnknown
		return
	}

	if maxSeq < conversation.StartSeq {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("max seq < start seq, max: %d, start seq: %d", maxSeq, conversation.StartSeq))
		err = code.ErrBadRequest
		return
	}

	if err = ConversationUseCase.UpdateCleanSeq(operationID, userID, conversationType, conversationID, maxSeq); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update clean seq error, conversation type: %d, conversation id: %s, error: %v", conversationType, conversationID, err))
		err = code.ErrDB
		return
	}
	return
}
