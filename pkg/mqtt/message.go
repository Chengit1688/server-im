package mqtt

import (
	"fmt"
	"im/pkg/common"
	"im/pkg/logger"
	"im/pkg/util"
)

// 发送消息到指定用户
func SendMessageToUsers(operationID string, mt common.MessageType, message interface{}, userIDList ...string) (err error) {
	for _, userID := range userIDList {
		topic := GetUserTopic(userID)
		if err = SendMessage(operationID, mt, message, topic); err != nil {
			return
		}
	}
	return
}

// 发送消息到指定群
func SendMessageToGroups(operationID string, mt common.MessageType, message interface{}, groupIDList ...string) (err error) {
	for _, groupID := range groupIDList {
		topic := GetGroupTopic(groupID)
		if err = SendMessage(operationID, mt, message, topic); err != nil {
			return
		}
	}
	return
}

// 广播消息到所有用户
func BroadcastMessage(operationID string, mt common.MessageType, message interface{}) (err error) {
	topic := GetSystemTopic()
	return SendMessage(operationID, mt, message, topic)
}

func SendMessage(operationID string, mt common.MessageType, message interface{}, topic string) (err error) {
	var msg common.Message
	msg.Type = mt
	msg.Data = message

	var data []byte
	if data, err = util.JsonMarshal(&msg); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("json marshal error, error: %v", err))
		return
	}

	if defaultEMQXClientManager == nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", "default emqx client manager nil")
		return
	}

	go defaultEMQXClientManager.Publish(operationID, topic, string(data))
	return
}
