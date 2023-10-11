package job

import (
	"fmt"
	"im/config"
	apiUserModel "im/internal/api/user/model"
	apiUserUseCase "im/internal/api/user/usecase"
	userUseCase "im/internal/cms_api/user/usecase"
	"im/pkg/common"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Init() {
	cfg := config.Config.EMQXServer

	client := mqtt.NewEMQXClient(cfg.MQTTAddress, cfg.MQTTUsername, cfg.MQTTPassword)
	go client.Subscribe(Callback, []string{mqtt.ClientConnectedTopic, mqtt.ClientDisconnectedTopic})
}

func Callback(client MQTT.Client, msg MQTT.Message) {
	logger.Sugar.Debugw("", "func", util.GetSelfFuncName(), "info", fmt.Sprintf("msg topic: %s", msg.Topic()))

	switch {

	case mqtt.ClientConnectedTopicRx.MatchString(msg.Topic()) || mqtt.ClientDisconnectedTopicRx.MatchString(msg.Topic()):
		cfg := config.Config
		var data mqtt.ClientConnPayload
		util.JsonUnmarshal(msg.Payload(), &data)

		if strings.Contains(data.ClientID, "_system_") {
			return
		}

		if data.Username == "root" {
			return
		}

		userID := getUserID(data.Username)
		siteName := getSiteName(data.Username, userID)
		if cfg.Station != siteName {

			logger.Sugar.Debugw("UpdateUserCallback", "topic", msg.Topic(), "siteName", siteName, "ignore", "其他站点数据 忽略")
			return
		}

		if !apiUserUseCase.UserUseCase.IsOnline(userID) {
			userUseCase.UserUseCase.BroadcastUserOnlineStatusToFriends(userID, common.UserStatusOffline)
			_ = userUseCase.UserUseCase.UpdateOnlineStatus(userID, apiUserModel.OnlineStatusTypeOffline, data.Ts/1000)
		} else {
			userUseCase.UserUseCase.BroadcastUserOnlineStatusToFriends(userID, common.UserStatusOnline)
			_ = userUseCase.UserUseCase.UpdateOnlineStatus(userID, apiUserModel.OnlineStatusTypeOnline, data.Ts/1000)
		}
	}
}

func getUserID(username string) (user_id string) {
	dataArray := strings.Split(username, "_")
	user_id = dataArray[len(dataArray)-1]
	return
}

func getSiteName(username, user_id string) (site_name string) {

	siteArray := strings.Split(username, user_id)
	site_name = siteArray[0][:len(siteArray[0])-1]

	return
}
