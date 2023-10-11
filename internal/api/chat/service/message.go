package service

import (
	"fmt"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	chatUseCase "im/internal/api/chat/usecase"
	friendUseCase "im/internal/api/friend/usecase"
	groupModel "im/internal/api/group/model"
	groupRepo "im/internal/api/group/repo"
	groupUseCase "im/internal/api/group/usecase"
	"im/internal/api/moments/model"
	permissionUseCase "im/internal/api/permission/usecase"
	"im/pkg/response"

	userUseCase "im/internal/api/user/usecase"
	"im/pkg/common"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

var MessageService = new(messageService)

type messageService struct{}

func (s *messageService) Pull(c *gin.Context) {
	var (
		req          chatModel.MessagePullReq
		resp         chatModel.MessagePullResp
		conversation *chatModel.Conversation
		err          error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.RecvID == "" {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "recv id nil")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	conversationID := chatUseCase.ConversationUseCase.GetConversationID(req.ConversationType, userID, req.RecvID)
	if conversation, err = chatRepo.ConversationRepo.Get(userID, req.ConversationType, conversationID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get error, conversation type: %d, conversation id: %s, error: %v", req.ConversationType, conversationID, err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	if resp.List, err = chatRepo.MessageRepo.PullNew(req.ConversationType, conversationID, conversation.StartSeq, req.Seq, req.PageSize); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db list error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	switch req.ConversationType {
	case chatModel.ConversationTypeSingle:
		user, err2 := userUseCase.UserUseCase.GetBaseInfo(req.RecvID)
		if err2 != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("single get base info error, user id: %s, error: %v", req.RecvID, err2))
		}
		userSelf, err3 := userUseCase.UserUseCase.GetBaseInfo(userID)
		if err3 != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("single get base info error, user id: %s, error: %v", userID, err3))
		}
		for i, message := range resp.List {

			if message.SendID != "" && message.SendID != userID && user != nil {
				message.SendNickname = user.NickName
				message.SendFaceUrl = user.FaceURL
			}

			if message.SendID != "" && message.SendID == userID && userSelf != nil {
				message.SendNickname = userSelf.NickName
				message.SendFaceUrl = userSelf.FaceURL
			}

			chatUseCase.MessageUseCase.FillSystemMessageUserInfo(req.OperationID, &message)

			message.RecvID = req.RecvID
			message.Content, _ = util.Encrypt([]byte(message.Content), common.ContentKey)
			resp.List[i] = message
		}

	case chatModel.ConversationTypeGroup:
		for i, message := range resp.List {

			if message.SendID != "" && message.SendID != userID {
				user, err2 := userUseCase.UserUseCase.GetBaseInfo(message.SendID)
				if err2 != nil {
					logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, user id: %s, error: %v", req.RecvID, err2))
				}

				if user != nil {
					message.SendNickname = user.NickName
					message.SendFaceUrl = user.FaceURL
				}
			}
			if message.SendID != "" && message.SendID == userID {
				userSelf, err3 := userUseCase.UserUseCase.GetBaseInfo(userID)
				if err3 != nil {
					logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, user id: %s, error: %v", userID, err3))
				}

				if userSelf != nil {
					message.SendNickname = userSelf.NickName
					message.SendFaceUrl = userSelf.FaceURL
				}
			}

			chatUseCase.MessageUseCase.FillSystemMessageUserInfo(req.OperationID, &message)

			message.RecvID = req.RecvID
			message.Content, _ = util.Encrypt([]byte(message.Content), common.ContentKey)
			resp.List[i] = message
		}
	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", req.ConversationType))
	}

	http.Success(c, resp)
	return
}

func (s *messageService) PullV2(c *gin.Context) {
	var (
		req          chatModel.MessagePullReq
		resp         chatModel.MessagePullResp
		conversation *chatModel.Conversation
		err          error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.RecvID == "" {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "recv id nil")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	conversationID := chatUseCase.ConversationUseCase.GetConversationID(req.ConversationType, userID, req.RecvID)
	if conversation, err = chatRepo.ConversationRepo.Get(userID, req.ConversationType, conversationID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get error, conversation type: %d, conversation id: %s, error: %v", req.ConversationType, conversationID, err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	if len(req.SeqList) > 0 {
		if resp.List, err = chatRepo.MessageRepo.PullNewList(req.ConversationType, conversationID, req.SeqList, req.PageSize); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db list error, error: %v", err))
			http.Failed(c, response.GetError(response.ErrDB, lang))
			return
		}
	} else {
		if resp.List, err = chatRepo.MessageRepo.PullNewV2(req.ConversationType, conversationID, conversation.StartSeq, req.Seq, req.PageSize); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db list error, error: %v", err))
			http.Failed(c, response.GetError(response.ErrDB, lang))
			return
		}
	}

	switch req.ConversationType {
	case chatModel.ConversationTypeSingle:
		user, err2 := userUseCase.UserUseCase.GetBaseInfo(req.RecvID)
		if err2 != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("single get base info error, user id: %s, error: %v", req.RecvID, err2))
		}
		userSelf, err3 := userUseCase.UserUseCase.GetBaseInfo(userID)
		if err3 != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("single get base info error, user id: %s, error: %v", userID, err3))
		}
		for i, message := range resp.List {

			if message.SendID != "" && message.SendID != userID && user != nil {
				message.SendNickname = user.NickName
				message.SendFaceUrl = user.FaceURL
			}

			if message.SendID != "" && message.SendID == userID && userSelf != nil {
				message.SendNickname = userSelf.NickName
				message.SendFaceUrl = userSelf.FaceURL
			}

			chatUseCase.MessageUseCase.FillSystemMessageUserInfo(req.OperationID, &message)

			message.RecvID = req.RecvID
			message.Content, _ = util.Encrypt([]byte(message.Content), common.ContentKey)
			resp.List[i] = message
		}

	case chatModel.ConversationTypeGroup:
		for i, message := range resp.List {

			if message.SendID != "" && message.SendID != userID {
				user, err2 := userUseCase.UserUseCase.GetBaseInfo(message.SendID)
				if err2 != nil {
					logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, user id: %s, error: %v", req.RecvID, err2))
				}

				if user != nil {
					message.SendNickname = user.NickName
					message.SendFaceUrl = user.FaceURL
				}
			}
			if message.SendID != "" && message.SendID == userID {
				userSelf, err3 := userUseCase.UserUseCase.GetBaseInfo(userID)
				if err3 != nil {
					logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get base info error, user id: %s, error: %v", userID, err3))
				}

				if userSelf != nil {
					message.SendNickname = userSelf.NickName
					message.SendFaceUrl = userSelf.FaceURL
				}
			}

			chatUseCase.MessageUseCase.FillSystemMessageUserInfo(req.OperationID, &message)

			message.RecvID = req.RecvID
			message.Content, _ = util.Encrypt([]byte(message.Content), common.ContentKey)
			resp.List[i] = message
		}
	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", req.ConversationType))
	}

	http.Success(c, resp)
	return
}

func (s *messageService) Send(c *gin.Context) {
	var (
		req     chatModel.MessageSendReq
		resp    chatModel.MessageSendResp
		err     error
		message *chatModel.MessageInfo
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.RecvID == "" {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "recv id nil")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")

	conversationID := chatUseCase.ConversationUseCase.GetConversationID(req.ConversationType, userID, req.RecvID)
	switch req.ConversationType {
	case chatModel.ConversationTypeSingle:
		if err = s.CheckFriendPermission(req.OperationID, userID, req.RecvID, lang); err != nil {
			http.Failed(c, err)
			return
		}

		if err = s.CheckFriendBlack(req.OperationID, req.RecvID, userID, lang); err != nil {
			friendUseCase.FriendUseCase.BlackFriendPush(req.OperationID, userID, req.RecvID)
			http.Failed(c, response.GetError(response.ErrFriendInBlack, lang))
			return
		}

		if message, err = chatUseCase.MessageUseCase.SendMessageToUsers(req.OperationID, req.ClientMsgID, req.ConversationType, conversationID, userID, req.Type, req.Content, userID, req.RecvID); err != nil {
			http.Failed(c, err)
			return
		}

	case chatModel.ConversationTypeGroup:
		if err = permissionUseCase.PermissionUseCase.CheckChatGroupPermission(req.OperationID, req.RecvID, userID, lang); err != nil {
			http.Failed(c, err)
			return
		}

		if message, err = chatUseCase.MessageUseCase.SendMessageToGroup(req.OperationID, req.ClientMsgID, req.ConversationType, conversationID, userID, req.Type, req.Content); err != nil {
			http.Failed(c, err)
			return
		}

		if len(req.AtList) > 0 {
			if err = s.GroupAtPush(req, conversationID, userID); err != nil {
				http.Failed(c, err)
				return
			}
		}

	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", req.ConversationType))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	util.CopyStructFields(&resp.MessageInfo, message)
	http.Success(c, resp)
	return
}

func (s *messageService) GroupAtPush(req chatModel.MessageSendReq, groupID string, userID string) (err error) {
	userInfo, err := userUseCase.UserUseCase.GetBaseInfo(userID)
	if err != nil {
		return
	}
	pushData := model.GroupAtPush{
		ConversationID:  groupID,
		Timestamp:       time.Now().Unix(),
		PublisherUserID: userID,
		FaceURL:         userInfo.FaceURL,
		Type:            int64(req.Type),
	}
	groupMember, err1 := groupUseCase.GroupMemberUseCase.GetMember(groupID, userID)
	pushData.PublisherNickname = userInfo.NickName
	if err1 == nil && groupMember != nil {
		logger.Sugar.Warnw(req.OperationID, "func", util.GetSelfFuncName(), "error", err1)
		pushData.PublisherNickname = groupMember.GroupNickName
	}

	if req.AtList[0] == "all" {
		req.AtList = groupRepo.GroupRepo.GroupMemberIdList(groupID)
	}
	for _, atUser := range req.AtList {
		pushData.FriendUserID = atUser
		if err2 := mqtt.SendMessageToUsers(req.OperationID, common.GroupMemberAtPush, pushData, atUser); err2 != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err2)
			continue
		}
	}
	return
}

func (s *messageService) MultiSend(c *gin.Context) {
	var (
		req  chatModel.MessageMultiSendReq
		resp chatModel.MessageMultiSendResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")

	if err = permissionUseCase.PermissionUseCase.CheckChatMultiSendPermission(req.OperationID, userID, lang); err != nil {
		http.Failed(c, err)
		return
	}

	if len(req.RecvIDList) != 0 {
		s.multiSendV1(c, userID, &req, &resp)
	} else {
		s.multiSendV2(c, userID, &req, &resp)
	}
}

func (s *messageService) multiSendV1(c *gin.Context, userID string, req *chatModel.MessageMultiSendReq, resp *chatModel.MessageMultiSendResp) {
	var err error
	lang := c.GetHeader("Locale")
	switch req.ConversationType {
	case chatModel.ConversationTypeSingle:
		for i := 0; i < len(req.RecvIDList); i++ {

			if i != 0 && i%50 == 0 {
				time.Sleep(time.Millisecond * 100)
			}
			recvID := req.RecvIDList[i]

			if err = s.CheckFriendBlack(req.OperationID, recvID, userID, lang); err != nil {
				friendUseCase.FriendUseCase.BlackFriendPush(req.OperationID, userID, recvID)
				continue
			}
			conversationID := chatUseCase.ConversationUseCase.GetConversationID(req.ConversationType, userID, recvID)
			if err = s.CheckFriendPermission(req.OperationID, userID, recvID, lang); err != nil {
				continue
			}
			if err = s.CheckFriendBlack(req.OperationID, recvID, userID, lang); err != nil {
				continue
			}

			if _, err = chatUseCase.MessageUseCase.SendMessageToUsers(req.OperationID, "", req.ConversationType, conversationID, userID, req.Type, req.Content, userID, recvID); err != nil {
				continue
			}

		}

	case chatModel.ConversationTypeGroup:
		for i := 0; i < len(req.RecvIDList); i++ {

			if i != 0 && i%50 == 0 {
				time.Sleep(time.Millisecond * 100)
			}

			recvID := req.RecvIDList[i]
			if err = permissionUseCase.PermissionUseCase.CheckChatGroupPermission(req.OperationID, recvID, userID, lang); err != nil {
				continue
			}

			if _, err = chatUseCase.MessageUseCase.SendMessageToGroup(req.OperationID, "", req.ConversationType, recvID, userID, req.Type, req.Content); err != nil {
				continue
			}

		}

	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", req.ConversationType))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	http.Success(c, resp)
	return
}

func (s *messageService) multiSendV2(c *gin.Context, userID string, req *chatModel.MessageMultiSendReq, resp *chatModel.MessageMultiSendResp) {
	var (
		message *chatModel.MessageInfo
		err     error
	)
	lang := c.GetHeader("Locale")
	switch req.ConversationType {
	case chatModel.ConversationTypeSingle:
		for i := 0; i < len(req.MessageList); i++ {
			clientMsg := req.MessageList[i]

			recvID := clientMsg.RecvID

			if err = s.CheckFriendBlack(req.OperationID, recvID, userID, lang); err != nil {
				friendUseCase.FriendUseCase.BlackFriendPush(req.OperationID, userID, recvID)
				continue
			}
			conversationID := chatUseCase.ConversationUseCase.GetConversationID(req.ConversationType, userID, recvID)
			if err = s.CheckFriendPermission(req.OperationID, userID, recvID, lang); err != nil {
				continue
			}
			if err = s.CheckFriendBlack(req.OperationID, recvID, userID, lang); err != nil {
				continue
			}

			if message, err = chatUseCase.MessageUseCase.SendMessageToUsers(req.OperationID, clientMsg.ClientMsgID, req.ConversationType, conversationID, userID, req.Type, req.Content, userID, recvID); err != nil {
				continue
			}

			resp.List = append(resp.List, *message)
		}

	case chatModel.ConversationTypeGroup:
		for i := 0; i < len(req.MessageList); i++ {
			clientMsg := req.MessageList[i]

			recvID := clientMsg.RecvID
			if err = permissionUseCase.PermissionUseCase.CheckChatGroupPermission(req.OperationID, recvID, userID, lang); err != nil {
				continue
			}

			if message, err = chatUseCase.MessageUseCase.SendMessageToGroup(req.OperationID, clientMsg.ClientMsgID, req.ConversationType, recvID, userID, req.Type, req.Content); err != nil {
				continue
			}

			resp.List = append(resp.List, *message)
		}

	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", req.ConversationType))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	http.Success(c, resp)
	return
}

func (s *messageService) Forward(c *gin.Context) {
	var (
		req     chatModel.MessageForwardReq
		resp    chatModel.MessageForwardResp
		message *chatModel.MessageInfo
		err     error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	for i := 0; i < len(req.MessageList); i++ {
		clientMsg := req.MessageList[i]
		conversationID := chatUseCase.ConversationUseCase.GetConversationID(clientMsg.ConversationType, userID, clientMsg.RecvID)

		switch clientMsg.ConversationType {
		case chatModel.ConversationTypeSingle:

			if err = s.CheckFriendBlack(req.OperationID, clientMsg.RecvID, userID, lang); err != nil {
				friendUseCase.FriendUseCase.BlackFriendPush(req.OperationID, userID, clientMsg.RecvID)
				continue
			}
			if err = s.CheckFriendPermission(req.OperationID, userID, clientMsg.RecvID, lang); err != nil {
				continue
			}
			if err = s.CheckFriendBlack(req.OperationID, clientMsg.RecvID, userID, lang); err != nil {
				continue
			}

			if message, err = chatUseCase.MessageUseCase.SendMessageToUsers(req.OperationID, clientMsg.ClientMsgID, clientMsg.ConversationType, conversationID, userID, clientMsg.Type, clientMsg.Content, userID, clientMsg.RecvID); err != nil {
				continue
			}

		case chatModel.ConversationTypeGroup:
			if err = permissionUseCase.PermissionUseCase.CheckChatGroupPermission(req.OperationID, clientMsg.RecvID, userID, lang); err != nil {
				continue
			}

			if message, err = chatUseCase.MessageUseCase.SendMessageToGroup(req.OperationID, clientMsg.ClientMsgID, clientMsg.ConversationType, clientMsg.RecvID, userID, clientMsg.Type, clientMsg.Content); err != nil {
				continue
			}

		default:
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", clientMsg.ConversationType))
			http.Failed(c, response.GetError(response.ErrBadRequest, lang))
			return
		}

		resp.List = append(resp.List, *message)
	}

	http.Success(c, resp)
	return
}

func (s *messageService) Change(c *gin.Context) {
	var (
		req         chatModel.MessageChangeReq
		resp        chatModel.MessageChangeResp
		message     *chatModel.MessageInfo
		messageType chatModel.MessageType
		err         error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")

	switch req.Status {
	case chatModel.MessageStatusTypeRevoke:
		messageType = chatModel.MessageRevoke

	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "message status type error")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if len(req.MsgIDList) == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "msg id list size 0")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	var msg *chatModel.Message
	if msg, err = chatRepo.MessageRepo.Get(req.MsgIDList[0]); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("db get error, msg id list: %v, error: %v", req.MsgIDList, err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	if msg.Type > chatModel.MessageOperation {

	}

	switch msg.ConversationType {
	case chatModel.ConversationTypeSingle:
		recvID := chatUseCase.ConversationUseCase.GetRecvID(userID, chatModel.ConversationTypeSingle, msg.ConversationID)
		if err = s.CheckFriendPermission(req.OperationID, userID, recvID, lang); err != nil {
			http.Failed(c, err)
			return
		}

		if err = permissionUseCase.PermissionUseCase.CheckChatRevokePermission(req.OperationID, userID, msg.SendID, msg.SendTime, lang); err != nil {
			http.Failed(c, err)
			return
		}

	case chatModel.ConversationTypeGroup:
		if err = permissionUseCase.PermissionUseCase.CheckChatGroupRevokePermission(req.OperationID, msg.ConversationID, userID, msg.SendID, msg.SendTime, lang); err != nil {
			http.Failed(c, err)
			return
		}
	}

	if message, err = chatUseCase.MessageUseCase.UpdateStatus(req.OperationID, userID, msg.ConversationType, msg.ConversationID, messageType, req.Status, req.MsgIDList); err != nil {
		http.Failed(c, err)
		return
	}

	if message != nil {
		util.CopyStructFields(&resp.MessageInfo, message)
	}

	http.Success(c, resp)
	return
}

func (s *messageService) Clear(c *gin.Context) {
	var (
		req  chatModel.MessageClearReq
		resp chatModel.MessageClearResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	if err = chatUseCase.MessageUseCase.ClearClientNew(req.OperationID, userID, req.ConversationType, req.ConversationID); err != nil {
		http.Failed(c, err)
		return
	}
	switch req.ConversationType {
	case 1:
		if err = chatUseCase.MessageUseCase.ClearClientNew(req.OperationID, userID, req.ConversationType, req.ConversationID); err != nil {
			http.Failed(c, err)
			return
		}
	default:
		if err = chatUseCase.MessageUseCase.ClearConversation(req.OperationID, userID, req.ConversationType, req.ConversationID, req.MaxSeq); err != nil {
			http.Failed(c, err)
			return
		}
	}

	util.CopyStructFields(&resp, &req)
	mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageClearPush, resp, userID)
	http.Success(c, resp)
	return
}

func (s *messageService) CheckFriendPermission(operationID string, userID string, friendID string, lang string) (err error) {

	if !friendUseCase.FriendUseCase.CheckFriend(userID, friendID) {
		return response.GetError(response.ErrFriendNotExist, lang)
	}
	return
}

func (s *messageService) CheckFriendBlack(operationID string, userID string, friendID string, lang string) (err error) {

	if !friendUseCase.FriendUseCase.CheckBlackFriend(userID, friendID) {
		return response.GetError(response.ErrFriendInBlack, lang)
	}
	return
}

func (s *messageService) CheckGroupPermission(operationID string, userID string, groupID string, lang string) (err error) {

	var group groupModel.GroupInfo
	if group, err = groupUseCase.GroupUseCase.GroupInfo(groupID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrDB, lang)
	}

	if group.Status == 2 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, error: %v", err))
		return response.GetError(response.ErrGroupNotExist, lang)
	}

	if !groupUseCase.GroupMemberUseCase.CheckMember(groupID, userID) {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check member error, error: %v", err))
		return response.GetError(response.ErrGroupNotMember, lang)
	}
	return
}
