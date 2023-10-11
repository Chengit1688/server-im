package usecase

import (
	"fmt"
	friendusecase "im/internal/api/friend/usecase"
	apiUserModel "im/internal/api/user/model"
	apiUserUseCase "im/internal/api/user/usecase"
	"im/internal/cms_api/user/model"
	userRepo "im/internal/cms_api/user/repo"
	"im/pkg/common"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"time"
)

var UserUseCase = new(userUseCase)

type userUseCase struct{}

func (c *userUseCase) UpdateOnlineStatus(userID string, status apiUserModel.OnlineStatusType, latestLoginTime int64) (err error) {
	return userRepo.UserRepo.UpdateOnlineStatus(userID, status, latestLoginTime)
}

func (c *userUseCase) CalibrationOnlineUsers() {
	users, err := userRepo.UserRepo.GetOneDayOnlineUsers()
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get one day online users error, error: %v", err))
		return
	}

	for _, user := range users {
		if apiUserUseCase.UserUseCase.IsOnline(user.UserID) {
			continue
		}

		userRepo.UserRepo.UpdateOnlineStatus(user.UserID, apiUserModel.OnlineStatusTypeOffline, 0)
		time.Sleep(time.Millisecond * 100)
	}
}

func (c *userUseCase) BroadcastUserOnlineStatusToFriends(userID string, messageType common.MessageType) {
	friends := friendusecase.FriendUseCase.GetUserFriendIdList(userID)
	msg := "好友上线"
	if messageType == common.UserStatusOffline {
		msg = "好友下线"
	}
	for _, f := range friends {

		if !apiUserUseCase.UserUseCase.IsOnline(f) {
			continue
		}
		online := model.UserOnlineStatus{
			UserID: userID,
			Type:   messageType,
			Msg:    msg,
		}
		_ = mqtt.SendMessageToUsers("user_online", messageType, online, f)
		time.Sleep(time.Millisecond * 50)
	}
}
