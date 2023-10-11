package usecase

import (
	"encoding/json"
	"fmt"
	chatModel "im/internal/api/chat/model"
	groupUseCase "im/internal/api/group/usecase"
	"im/internal/api/push/model"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/common"
	"im/pkg/logger"
	"im/pkg/push"
	"im/pkg/util"
	"math"
	"regexp"
)

var JpushUseCase = new(jpushUseCase)

type jpushUseCase struct{}

func (j *jpushUseCase) Push(operationID, senderID, recvID, msg string, msgType chatModel.MessageType, conversationType chatModel.ConversationType) {
	var title, alert string
	alert, err := j.formatMsg(msg, msgType)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("formatMsg error, error: %v", err), "msgType", msgType, "msg", msg)
		return
	}
	user, err := userUseCase.UserUseCase.GetBaseInfo(senderID)
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("userUseCase GetBaseInfo, error: %v", err), "user_id", senderID)
		return
	}

	switch conversationType {
	case chatModel.ConversationTypeSingle:

		title = fmt.Sprintf("好友%s发送新消息：", user.NickName)
		titleRune := []rune(title)
		if len(titleRune) >= 20 {
			title = string(titleRune[:16]) + "..."
		}
		push.Jpush([]string{recvID}, alert, title, operationID)
	case chatModel.ConversationTypeGroup:

		users := groupUseCase.GroupUseCase.GroupMemberIdList(recvID)

		users = deleteSlice(users, user.UserId)
		group, err := groupUseCase.GroupUseCase.GroupInfo(recvID)
		if err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GroupUseCase GroupInfo, error: %v", err), "group_id", recvID)
			return
		}
		title = fmt.Sprintf("群%s %s发送新消息：", group.Name, user.NickName)
		titleRune := []rune(title)
		if len(titleRune) >= 20 {
			title = string(titleRune[:16]) + "..."
		}
		j.limitPush(users, alert, title, operationID)
	}
	return
}

func (j *jpushUseCase) formatMsg(msg string, msgType chatModel.MessageType) (message string, err error) {
	var decryptContent string
	var textMsg model.MsgText
	switch msgType {
	case chatModel.MessageText:

		var emojiRx = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)
		if decryptContent, err = util.Decrypt(msg, common.ContentKey); err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("decrypt error, error: %v", err))
		}
		json.Unmarshal([]byte(decryptContent), &textMsg)
		message = emojiRx.ReplaceAllString(textMsg.Text, `[表情]`)
		messageRune := []rune(message)
		if len(messageRune) >= 100 {
			message = string(messageRune[:96]) + "..."
		}
	case chatModel.MessageFace:

		var emojiRx = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)
		if decryptContent, err = util.Decrypt(msg, common.ContentKey); err != nil {
			logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("decrypt error, error: %v", err))
		}
		json.Unmarshal([]byte(decryptContent), &textMsg)
		message = emojiRx.ReplaceAllString(textMsg.Text, `[表情]`)
		messageRune := []rune(message)
		if len(messageRune) >= 100 {
			message = string(messageRune[:96]) + "..."
		}
	case chatModel.MessageImage:

		message = "[图片]"
	case chatModel.MessageVoice:

		message = "[语音]"
	case chatModel.MessageVideo:

		message = "[视频]"
	case chatModel.MessageFile:

		message = "[文件]"
	default:
		message = "您收到一条新消息"
	}
	return message, err
}

func (j *jpushUseCase) limitPush(users []string, alert, title, operationID string) {

	size := 1000
	lens := len(users)
	if lens > size {
		mod := math.Ceil(float64(lens) / float64(size))
		spliltList := make([][]string, 0)
		for i := 0; i < int(mod); i++ {
			tmpList := make([]string, 0, size)
			fmt.Println("i=", i)
			if i == int(mod)-1 {
				tmpList = users[i*size:]
			} else {
				tmpList = users[i*size : i*size+size]
			}
			spliltList = append(spliltList, tmpList)
		}
		for _, subUsers := range spliltList {
			push.Jpush(subUsers, alert, title, operationID)
		}
	} else {
		push.Jpush(users, alert, title, operationID)
	}
}

func deleteSlice(a []string, elem string) []string {
	tgt := a[:0]
	for _, v := range a {
		if v != elem {
			tgt = append(tgt, v)
		}
	}
	return tgt
}
