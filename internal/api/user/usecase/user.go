package usecase

import (
	"context"
	"fmt"
	"im/config"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	"im/pkg/code"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"strings"
	"sync"
	"time"
)

var UserUseCase = new(userUseCase)

type userUseCase struct{}

func (c *userUseCase) GetAllUserIDList() (userIDList []string, err error) {
	return userRepo.UserRepo.ListUserID()
}

func (c *userUseCase) GetBaseInfo(userID string) (user *userModel.UserBaseInfo, err error) {
	userInfo, e := userRepo.UserCache.GetBaseUserInfo(userID, userRepo.UserRepo.GetBaseInfoByUserId)
	if userInfo == nil {
		return &userModel.UserBaseInfo{}, code.ErrUserIdNotExist
	}
	return userInfo, e
}

func (c *userUseCase) GetBaseInfoList(userIDList []string) (users []userModel.UserBaseInfo) {
	return
}

func (c *userUseCase) GetInfo(userID string) (user *userModel.User, err error) {
	return userRepo.UserRepo.GetByUserID(userRepo.WhereOption{UserId: userID})
}

func (c *userUseCase) GetInfoList(userIDList []string) (users []userModel.UserBaseInfo) {
	return
}

func (c *userUseCase) IsOnline(userID string) bool {
	station := config.Config.Station
	username := fmt.Sprintf("%s_%s", station, userID)
	clients, err := mqtt.GetClients(username, "", "", mqtt.ConnStateTypeConnected, 0, 0)
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get clients error, error: %v", err))
		return false
	}
	return len(clients) != 0
}

func (c *userUseCase) OnlineMap() (users map[string]bool, err error) {
	station := config.Config.Station
	users = make(map[string]bool)
	sep := fmt.Sprintf("%s_", station)

	var clients []mqtt.Client
	if clients, err = mqtt.GetClients("", station, "", mqtt.ConnStateTypeConnected, 0, 0); err != nil {
		return
	}
	for _, client := range clients {
		_, userID, _ := strings.Cut(client.Username, sep)
		users[userID] = true
	}
	return
}

func (c *userUseCase) GetClientIDList(userID string) (clientIDList []string, err error) {
	station := config.Config.Station
	username := fmt.Sprintf("%s_%s", station, userID)

	clients, err := mqtt.GetClients(username, "", "", "", 0, 0)
	if err != nil {
		return
	}

	for _, client := range clients {
		clientIDList = append(clientIDList, client.ClientID)
	}
	return
}

type FuncForUserId func(userId string) error

func (c *userUseCase) DoRoutineByUserId(operationId, userId string, handlers ...FuncForUserId) {
	var (
		withTimeout context.Context
		cancelFunc  context.CancelFunc
		waitGroup   sync.WaitGroup
	)
	withTimeout, cancelFunc = context.WithTimeout(context.Background(), time.Second*5)
	waitGroup.Add(len(handlers))
	for _, f := range handlers {

		go func(handle func(userId string) error) {
			defer func() {
				if e := recover(); e != nil {
					logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("DoRoutineByUserId,error: %v", e))
				}
				waitGroup.Done()
			}()
			if err := handle(userId); err != nil {
				logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("DoRoutineByUserId,error: %v", err))
			}
		}(f)
	}
	go func() {
		select {
		case <-withTimeout.Done():
			return
		default:
			waitGroup.Wait()
			cancelFunc()
			return
		}
	}()
	<-withTimeout.Done()
}

func (c *userUseCase) GetBaseInfoByPhoneNumber(phoneNumber string) (*userModel.UserBaseInfo, error) {
	return userRepo.UserRepo.GetBaseInfoByPhoneNUmber(phoneNumber)
}

func (c *userUseCase) UpdateUserWallet(userID string, flag string, amount int64) error {
	return userRepo.UserRepo.UpdateWallet(userID, flag, amount)
}

func (c *userUseCase) CheckUserOnline(userId string) bool {
	return userRepo.UserRepo.CheckUserOnline(userId)
}
