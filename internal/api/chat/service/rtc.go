package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	chatUseCase "im/internal/api/chat/usecase"
	friendUseCase "im/internal/api/friend/usecase"
	permissionUseCase "im/internal/api/permission/usecase"
	"im/pkg/common"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"
	"time"
)

var RTCService = new(rtcService)

type rtcService struct{}

func (s *rtcService) RTCInfo(c *gin.Context) {
	var (
		req  chatModel.RTCInfoReq
		resp chatModel.RTCInfoResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	resp.RTCInfo, _ = chatRepo.RTCCache.GetRTC(userID)
	http.Success(c, resp)
	return
}

func (s *rtcService) RTC(c *gin.Context) {
	var (
		req      chatModel.RTCReq
		resp     chatModel.RTCResp
		sendInfo *chatModel.RTCInfo
		recvInfo *chatModel.RTCInfo
		err      error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.RTCType != chatModel.RTCTypeAudio && req.RTCType != chatModel.RTCTypeVideo {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc type error, rtc type: %d", req.RTCType))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if err = permissionUseCase.PermissionUseCase.CheckChatRTCPermission(req.OperationID, req.RTCType, lang); err != nil {
		http.Failed(c, err)
		return
	}

	userID := c.GetString("user_id")

	if sendInfo, err = chatRepo.RTCCache.GetRTC(userID); err != nil && err != redis.Nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get send rtc error, user id: %s, error: %v", userID, err))
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	if sendInfo != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send info exist, user id: %s, error: %v", userID, err))
		http.Failed(c, response.GetError(response.ErrChatRTCBusy, lang))
		return
	}

	switch req.ConversationType {
	case chatModel.ConversationTypeSingle:
		if err = MessageService.CheckFriendPermission(req.OperationID, userID, req.RecvID, lang); err != nil {
			http.Failed(c, err)
			return
		}

		if !friendUseCase.FriendUseCase.CheckBlackFriend(req.RecvID, userID) {
			friendUseCase.FriendUseCase.BlackFriendPush(req.OperationID, userID, req.RecvID)
			http.Failed(c, response.GetError(response.ErrFriendInBlack, lang))
			return
		}

		if recvInfo, err = chatRepo.RTCCache.GetRTC(req.RecvID); err != nil && err != redis.Nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get recv rtc error, user id: %s, error: %v", req.RecvID, err))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if recvInfo != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("recv info exist, user id: %s, error: %v", req.RecvID, err))
			http.Failed(c, response.GetError(response.ErrChatRTCTargetBusy, lang))

			recvInfo.SendID = userID
			recvInfo.RecvID = req.RecvID
			chatUseCase.RTCUseCase.Finish(req.OperationID, recvInfo, chatModel.RTCStatusTypeBusy)
			return
		}

		sendInfo = chatUseCase.RTCUseCase.FillRequestInfo(req.ConversationType, userID, req.RecvID, req.RTCType, req.DeviceID)

		redisExpireTime := chatModel.RTCRequestExpireTime + 10
		_ = chatRepo.RTCCache.SetRTC(userID, sendInfo, redisExpireTime)
		_ = chatRepo.RTCCache.SetRTC(req.RecvID, sendInfo, redisExpireTime)

		_ = mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageRTCPush, sendInfo, userID, req.RecvID)

	case chatModel.ConversationTypeGroup:

		if err = MessageService.CheckGroupPermission(req.OperationID, userID, req.RecvID, lang); err != nil {
			http.Failed(c, err)
			return
		}
	}

	resp.RTCInfo = sendInfo
	http.Success(c, resp)
	return
}

func (s *rtcService) RTCOperate(c *gin.Context) {
	var (
		req     chatModel.RTCOperateReq
		resp    chatModel.RTCOperateResp
		rtcInfo *chatModel.RTCInfo
		err     error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "info", req)

	userID := c.GetString("user_id")

	if rtcInfo, err = chatRepo.RTCCache.GetRTC(userID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get rtc error, user id: %s, error: %v", userID, err))
		http.Failed(c, response.GetError(response.ErrChatRTCNotFound, lang))
		return
	}

	switch req.OperationType {
	case chatModel.RTCOperationTypeCancel:

		if rtcInfo.SendID != userID {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("cancel error, send id: %s, user id: %s", rtcInfo.SendID, userID))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeRequest {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if err = chatUseCase.RTCUseCase.CheckUserAndDevice(userID, req.DeviceID, rtcInfo); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user and device error, user id: %s", userID))
			http.Failed(c, err)
			return
		}

		rtcInfo.SendDeviceID = req.DeviceID
		rtcInfo = chatUseCase.RTCUseCase.Finish(req.OperationID, rtcInfo, chatModel.RTCStatusTypeCancel)

	case chatModel.RTCOperationTypeAgree:

		if rtcInfo.RecvID != userID {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("agree error, recv id: %s, user id: %s", rtcInfo.RecvID, userID))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeRequest {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		now := time.Now().Unix()
		redisExpireTime := chatModel.RTCRetainExpireTime
		rtcInfo.RecvDeviceID = req.DeviceID
		rtcInfo.RTCStatus = chatModel.RTCStatusTypeAgree
		rtcInfo.RTCStartTime = now
		rtcInfo.RTCUpdateTime = now
		_ = chatRepo.RTCCache.SetRTC(rtcInfo.SendID, rtcInfo, redisExpireTime)
		_ = chatRepo.RTCCache.SetRTC(rtcInfo.RecvID, rtcInfo, redisExpireTime)

		_ = mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageRTCPush, rtcInfo, rtcInfo.SendID, rtcInfo.RecvID)

	case chatModel.RTCOperationTypeDisagree:

		if rtcInfo.RecvID != userID {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("disagree error, recv id: %s, user id: %s", rtcInfo.RecvID, userID))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeRequest {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}
		tmpID := rtcInfo.SendID
		rtcInfo.SendID = userID
		rtcInfo.RecvID = tmpID
		rtcInfo.RecvDeviceID = req.DeviceID
		logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc rtcInfo: %+v", rtcInfo))
		rtcInfo = chatUseCase.RTCUseCase.Finish(req.OperationID, rtcInfo, chatModel.RTCStatusTypeDisagree)

	case chatModel.RTCOperationTypeFinish:

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeAgree {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if err = chatUseCase.RTCUseCase.CheckUserAndDevice(userID, req.DeviceID, rtcInfo); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user and device error, user id: %s", userID))
			http.Failed(c, err)
			return
		}

		rtcInfo = chatUseCase.RTCUseCase.Finish(req.OperationID, rtcInfo, chatModel.RTCStatusTypeFinish)

	case chatModel.RTCOperationTypeSwitch:
		if req.RTCType != chatModel.RTCTypeAudio && req.RTCType != chatModel.RTCTypeVideo {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc type error, rtc type: %d", req.RTCType))
			http.Failed(c, response.GetError(response.ErrBadRequest, lang))
			return
		}

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeRequest && rtcInfo.RTCStatus != chatModel.RTCStatusTypeAgree {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if err = chatUseCase.RTCUseCase.CheckUserAndDevice(userID, req.DeviceID, rtcInfo); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user and device error, user id: %s", userID))
			http.Failed(c, err)
			return
		}

		redisExpireTime := chatModel.RTCRequestExpireTime
		rtcInfo.RTCType = req.RTCType
		_ = chatRepo.RTCCache.SetRTC(rtcInfo.SendID, rtcInfo, redisExpireTime)
		_ = chatRepo.RTCCache.SetRTC(rtcInfo.RecvID, rtcInfo, redisExpireTime)

		_ = mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageRTCPush, rtcInfo, rtcInfo.SendID, rtcInfo.RecvID)

	case chatModel.RTCOperationTypeNotResponse:

		if rtcInfo.SendID != userID {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("cancel error, send id: %s, user id: %s", rtcInfo.SendID, userID))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeRequest {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if err = chatUseCase.RTCUseCase.CheckUserAndDevice(userID, req.DeviceID, rtcInfo); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user and device error, user id: %s", userID))
			http.Failed(c, err)
			return
		}

		rtcInfo.SendDeviceID = req.DeviceID
		logger.Sugar.Errorf("请求参数:%+v,发送取消消息状态RTCStatusTypeNotResponse:%d", req, chatModel.RTCStatusTypeNotResponse)
		rtcInfo = chatUseCase.RTCUseCase.Finish(req.OperationID, rtcInfo, chatModel.RTCStatusTypeNotResponse)

	case chatModel.RTCOperationTypeAbort:

		if rtcInfo.RTCStatus != chatModel.RTCStatusTypeAgree {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}

		if err = chatUseCase.RTCUseCase.CheckUserAndDevice(userID, req.DeviceID, rtcInfo); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user and device error, user id: %s", userID))
			http.Failed(c, err)
			return
		}

		rtcInfo = chatUseCase.RTCUseCase.Finish(req.OperationID, rtcInfo, chatModel.RTCStatusTypeAbort)

	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bad request, operation type: %d", req.OperationType))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	resp.RTCInfo = rtcInfo
	http.Success(c, resp)
	return
}

func (s *rtcService) RTCUpdate(c *gin.Context) {
	var (
		req  chatModel.RTCUpdateReq
		resp chatModel.RTCUpdateResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")

	rtcInfo, _ := chatRepo.RTCCache.GetRTC(userID)
	if rtcInfo == nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get rtc error, user id: %s", userID))
		http.Failed(c, response.GetError(response.ErrChatRTCNotFound, lang))
		return
	}

	if rtcInfo.RTCStatus != chatModel.RTCStatusTypeAgree {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
		http.Failed(c, response.GetError(response.ErrChatRTCStatus, lang))
		return
	}

	if err = chatUseCase.RTCUseCase.CheckUserAndDevice(userID, req.DeviceID, rtcInfo); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("check user and device error, user id: %s", userID))
		http.Failed(c, err)
		return
	}

	now := time.Now().Unix()
	rtcInfo.RTCUpdateTime = now
	_ = chatRepo.RTCCache.SetRTC(userID, rtcInfo, chatModel.RTCRetainExpireTime)

	if rtcInfo.SendID == userID {
		if err = chatUseCase.RTCUseCase.CheckAgreeStatus(req.OperationID, rtcInfo, rtcInfo.RecvID); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("receiver check agree status error, user id: %s", userID))
			http.Failed(c, err)
			return
		}
	} else if rtcInfo.RecvID == userID {
		if err = chatUseCase.RTCUseCase.CheckAgreeStatus(req.OperationID, rtcInfo, rtcInfo.SendID); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("sender check agree status error, user id: %s", userID))
			http.Failed(c, err)
			return
		}
	} else {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "user id not sender or receiver")
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	http.Success(c, resp)
	return
}
