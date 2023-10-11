package usecase

import (
	"context"
	"fmt"
	"im/config"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
)

var RTCUseCase = new(rtcUseCase)

type rtcUseCase struct{}

func (c *rtcUseCase) FillRequestInfo(conversationType chatModel.ConversationType, sendID string, recvID string, rtcType chatModel.RTCType, deviceID string) (rtcInfo *chatModel.RTCInfo) {
	conversationID := ConversationUseCase.GetConversationID(conversationType, sendID, recvID)

	rtcInfo = new(chatModel.RTCInfo)
	rtcInfo.ConversationType = conversationType
	rtcInfo.SendID = sendID
	rtcInfo.RecvID = recvID
	rtcInfo.RTCChannel, rtcInfo.RTCToken, _ = c.GenerateToken(conversationType, conversationID)
	rtcInfo.RTCType = rtcType
	rtcInfo.RTCStatus = chatModel.RTCStatusTypeRequest
	rtcInfo.RTCStartTime = time.Now().Unix()
	rtcInfo.RTCRequestLimitTime = chatModel.RTCRequestExpireTime
	rtcInfo.SendDeviceID = deviceID

	sendUser, err := userUseCase.UserUseCase.GetBaseInfo(sendID)
	if err != nil {
		return
	}
	rtcInfo.SendNickname = sendUser.NickName
	rtcInfo.SendFaceURL = sendUser.FaceURL

	recvUser, err := userUseCase.UserUseCase.GetBaseInfo(recvID)
	if err != nil {
		return
	}
	rtcInfo.RecvNickname = recvUser.NickName
	rtcInfo.RecvFaceURL = recvUser.FaceURL
	return
}

func (c *rtcUseCase) CheckUserAndDevice(userID string, deviceID string, rtcInfo *chatModel.RTCInfo) (err error) {
	if rtcInfo.SendID == userID {
		if rtcInfo.SendDeviceID != deviceID {
			err = code.ErrChatRTCDevice
		}
	} else if rtcInfo.RecvID == userID {
		if rtcInfo.RecvDeviceID != deviceID {
			err = code.ErrChatRTCDevice
		}
	} else {
		err = code.ErrChatRTCUser
	}
	return
}

func (c *rtcUseCase) CheckAgreeStatus(operationID string, info *chatModel.RTCInfo, userID string) (err error) {

	rtcInfo, _ := chatRepo.RTCCache.GetRTC(userID)
	if rtcInfo == nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc info nil, user id: %s", userID))
		err = code.ErrChatRTCAbort
		return
	}

	if rtcInfo.RTCStatus != chatModel.RTCStatusTypeAgree {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc status error, user id: %s, rtc status: %d", userID, rtcInfo.RTCStatus))
		err = code.ErrChatRTCStatus
		return
	}

	now := time.Now().Unix()
	if now-rtcInfo.RTCUpdateTime >= 5 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc finish, user id: %s", userID))
		err = code.ErrChatRTCAbort
	} else if now-rtcInfo.RTCUpdateTime >= 3 {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("rtc network bad, user id: %s", userID))
		err = code.ErrChatRTCNetworkBad
	}

	if err == code.ErrChatRTCAbort {
		c.Finish(operationID, info, chatModel.RTCStatusTypeAbort)
	}
	return
}

func (c *rtcUseCase) StartRequestTask(operationID string, sendID string) {
	go func() {
		var (
			sendInfo *chatModel.RTCInfo
			recvInfo *chatModel.RTCInfo
			err      error
		)

		expireTime := time.Duration(chatModel.RTCRequestExpireTime + 5)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*expireTime)
		defer cancel()

		select {
		case <-ctx.Done():
			if sendInfo, err = chatRepo.RTCCache.GetRTC(sendID); err != nil {
				return
			}

			switch sendInfo.ConversationType {
			case chatModel.ConversationTypeSingle:
				if recvInfo, err = chatRepo.RTCCache.GetRTC(sendInfo.RecvID); err != nil {
					return
				}
				logger.Sugar.Errorf("请求状态send:%d,请求状态recv:%d", chatModel.RTCStatusTypeRequest, chatModel.RTCStatusTypeRequest)

				if sendInfo.RTCStatus == chatModel.RTCStatusTypeRequest && recvInfo.RTCStatus == chatModel.RTCStatusTypeRequest {
					logger.Sugar.Errorf("发送取消消息状态:%d", chatModel.RTCStatusTypeNotResponse)
					c.Finish(operationID, sendInfo, chatModel.RTCStatusTypeNotResponse)
				}
			case chatModel.ConversationTypeGroup:

			}
		}
	}()
}

func (c *rtcUseCase) Finish(operationID string, rtcInfo *chatModel.RTCInfo, rtcStatus chatModel.RTCStatusType) *chatModel.RTCInfo {

	if rtcStatus != chatModel.RTCStatusTypeBusy {
		chatRepo.RTCCache.DeleteRTC(chatRepo.RTCCache.GetRTCKey(rtcInfo.SendID), chatRepo.RTCCache.GetRTCKey(rtcInfo.RecvID))
	}

	var data chatModel.MessageContent
	data.OperatorID = rtcInfo.SendID

	content := struct {
		SendID      string                  `json:"send_id"`
		RecvID      string                  `json:"recv_id"`
		RTCType     chatModel.RTCType       `json:"rtc_type"`
		RTCStatus   chatModel.RTCStatusType `json:"rtc_status"`
		RTCDuration int64                   `json:"rtc_duration"`
	}{}

	content.SendID = rtcInfo.SendID
	content.RecvID = rtcInfo.RecvID
	content.RTCType = rtcInfo.RTCType
	content.RTCStatus = rtcStatus

	if rtcStatus == chatModel.RTCStatusTypeFinish || rtcStatus == chatModel.RTCStatusTypeAbort {
		content.RTCDuration = time.Now().Unix() - rtcInfo.RTCStartTime
	}

	byteData, err := util.JsonMarshal(&content)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("json marshal error, error: %v", err))
		return nil
	}

	data.Content = string(byteData)

	conversationID := ConversationUseCase.GetConversationID(rtcInfo.ConversationType, rtcInfo.SendID, rtcInfo.RecvID)
	_, _ = MessageUseCase.SendSystemMessageToUsers(operationID, rtcInfo.ConversationType, conversationID, chatModel.MessageRTCNotify, &data, rtcInfo.SendID, rtcInfo.RecvID)

	rtcInfo.RTCStatus = rtcStatus
	if rtcStatus == chatModel.RTCStatusTypeBusy {

		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessageRTCPush, rtcInfo, rtcInfo.SendID)
	} else {

		_ = mqtt.SendMessageToUsers(operationID, common.ChatMessageRTCPush, rtcInfo, rtcInfo.SendID, rtcInfo.RecvID)
	}
	return rtcInfo
}

func (c *rtcUseCase) GenerateToken(conversationType chatModel.ConversationType, conversationID string) (rtcChannel string, rtcToken string, err error) {
	cfg := config.Config.Agora

	appID := cfg.AppID
	appSecret := cfg.AppSecret
	tokenExpireTime := cfg.TokenExpireTime

	if tokenExpireTime == 0 {
		tokenExpireTime = 2
	}
	expireTimestamp := time.Now().UTC().Unix() + tokenExpireTime*3600
	channelName := c.generateChannelName(conversationType, conversationID)

	var role rtctokenbuilder.Role = rtctokenbuilder.RolePublisher
	rtcToken, err = rtctokenbuilder.BuildTokenWithUID(appID, appSecret, channelName, 0, role, uint32(expireTimestamp))
	return channelName, rtcToken, err
}

func (c *rtcUseCase) generateChannelName(conversationType chatModel.ConversationType, conversationID string) string {
	var name string
	switch conversationType {
	case chatModel.ConversationTypeSingle:
		name = "single"
	case chatModel.ConversationTypeGroup:
		name = "group"
	}
	return fmt.Sprintf("%s_%s_%s_%s", config.Config.Station, name, conversationID, util.RandID(10))
}
