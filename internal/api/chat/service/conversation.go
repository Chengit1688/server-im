package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	chatUseCase "im/internal/api/chat/usecase"
	groupUseCase "im/internal/api/group/usecase"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/common"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var ConversationService = new(conversationService)

type conversationService struct{}

func (s *conversationService) List(c *gin.Context) {
	var (
		req  chatModel.ConversationListReq
		resp chatModel.ConversationListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	req.Check()
	resp.Pagination = req.Pagination
	userID := c.GetString("user_id")

	if req.Version < 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "version < 0")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	list, count, err := chatRepo.ConversationRepo.List(userID, req.Version, req.Offset, req.Limit)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db list error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	for _, conversation := range list {
		var data chatModel.ConversationInfo
		data.Conversation = conversation

		seq := conversation.AckSeq
		if seq < conversation.StartSeq {
			seq = conversation.StartSeq
		}
		switch conversation.ConversationType {
		case chatModel.ConversationTypeSingle:
			recvID := chatUseCase.ConversationUseCase.GetRecvID(userID, chatModel.ConversationTypeSingle, conversation.ConversationID)
			user, err2 := userUseCase.UserUseCase.GetBaseInfo(recvID)
			if err2 != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("single get base info error, user id: %s, error: %v", recvID, err2))
			}

			if user != nil {
				data.ConversationName = user.NickName
				data.ConversationFaceUrl = user.FaceURL
			}

			data.Message.RecvID = chatUseCase.ConversationUseCase.GetRecvID(userID, conversation.ConversationType, conversation.ConversationID)

			if data.Message.SendID != "" && data.Message.SendID != userID {
				data.Message.SendNickname = data.ConversationName
				data.Message.SendFaceUrl = data.ConversationFaceUrl
			}
		case chatModel.ConversationTypeGroup:

			group, err2 := groupUseCase.GroupUseCase.GroupInfo(conversation.ConversationID)
			if err2 != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group info error, group id: %s, error: %v", conversation.ConversationID, err2))
			}

			if group.GroupId != "" {
				data.ConversationName = group.Name
				data.ConversationFaceUrl = group.FaceUrl
			}

			data.Message.RecvID = conversation.ConversationID

			if data.Message.SendID != "" && data.Message.SendID != userID {
				user, err3 := userUseCase.UserUseCase.GetBaseInfo(data.Message.SendID)
				if err3 != nil {
					logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("group get base info error, user id: %s, error: %v", user.UserId, err3))
				}
				data.Message.SendNickname = user.NickName
				data.Message.SendFaceUrl = user.FaceURL
			}
		default:
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation type not found, conversation type: %d", conversation.ConversationType))
		}

		chatUseCase.MessageUseCase.FillSystemMessageUserInfo(req.OperationID, data.Message)

		data.Message.Content, err = util.Encrypt([]byte(data.Message.Content), common.ContentKey)

		if data.Message.ID == 0 {
			data.Message = nil
		}

		resp.List = append(resp.List, data)
	}
	resp.Count = count
	http.Success(c, resp)
	return
}

func (s *conversationService) AckSeq(c *gin.Context) {
	var (
		req  chatModel.ConversationAckSeqReq
		resp chatModel.ConversationAckSeqResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	var conversation *chatModel.Conversation
	if conversation, err = chatRepo.ConversationRepo.Get(userID, req.ConversationType, req.ConversationID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	if conversation.ID == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("conversation not found, conversation type: %d, conversation id: %s", req.ConversationType, req.ConversationID))
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	if req.AckSeq < conversation.AckSeq {

		http.Success(c, resp)
		return
	}

	var maxSeq int64
	s.pushReadSeq(req.OperationID, req.ConversationType, req.ConversationID, maxSeq)

	if err = chatUseCase.ConversationUseCase.UpdateAckSeq(req.OperationID, userID, req.ConversationType, req.ConversationID, maxSeq); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update ack seq error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	util.CopyStructFields(&resp, &req)
	mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageAckSeqPush, resp, userID)
	http.Success(c, resp)
	return
}

func (s *conversationService) pushReadSeq(operationID string, conversationType chatModel.ConversationType, conversationID string, ackSeq int64) (err error) {
	var readSeq int64
	if readSeq, err = chatUseCase.ConversationUseCase.GetReadSeq(conversationType, conversationID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get read seq error, error: %v", err))
		return
	}

	if ackSeq <= readSeq {
		return
	}

	var pushMsg chatModel.ConversationReadSeqResp
	pushMsg.ConversationType = conversationType
	pushMsg.ConversationID = conversationID
	pushMsg.ReadSeq = ackSeq

	switch conversationType {
	case chatModel.ConversationTypeSingle:
		userIDList := strings.Split(conversationID, "_")
		if len(userIDList) != 2 {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("split conversation id error, size != 2, conversation id: %s", conversationID))
			return
		}

		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessageReadSeqPush, pushMsg, userIDList...)

	case chatModel.ConversationTypeGroup:
		_ = mqtt.SendMessageToGroups(operationID, common.ChatMessageReadSeqPush, pushMsg, conversationID)
	}

	return
}
