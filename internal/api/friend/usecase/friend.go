package usecase

import (
	"fmt"
	"im/internal/api/friend/model"
	"im/internal/api/friend/repo"
	userModel "im/internal/api/user/model"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

var FriendUseCase = new(friendUseCase)

type friendUseCase struct{}

func (c *friendUseCase) GetFriendInfo(userID string, friendID string) (friendInfo *model.FriendInfo, err error) {
	var friend *model.Friend
	if friend, err = c.GetFriend(userID, friendID); err != nil {
		return
	}

	user, err2 := userUseCase.UserUseCase.GetBaseInfo(friendID)
	if err2 != nil {
		logger.Sugar.Errorw("", "func", "GetFriend", "error", fmt.Sprintf("get user base info error, error: %v", err2))
		return
	}

	friendInfo = new(model.FriendInfo)
	util.CopyStructFields(friendInfo, friend)
	friendInfo.ID = friend.ID
	friendInfo.NickName = user.NickName
	friendInfo.FaceURL = user.FaceURL
	friendInfo.BigFaceURL = user.BigFaceURL
	friendInfo.Signatures = user.Signatures
	friendInfo.Gender = user.Gender
	friendInfo.UserId = friendID
	friendInfo.Account = user.Account
	friendInfo.Age = user.Age
	friendInfo.OnlineStatus = 1
	friendInfo.PhoneNumber = user.PhoneNumber
	return
}

func (c *friendUseCase) GetFriend(userID string, friendID string) (friend *model.Friend, err error) {
	friend, err = repo.FriendCache.GetFriend(userID, friendID)
	if err != nil && err != redis.Nil {
		return
	}

	if err == nil {
		return
	}

	key := repo.FriendCache.GetFriendKey(userID, friendID)
	l := util.NewLock(db.RedisCli, key)
	if err = l.Lock(); err != nil {
		return
	}
	defer l.Unlock()

	friend, err = repo.FriendCache.GetFriend(userID, friendID)
	if err == nil {
		return
	}

	if friend, err = repo.FriendRepo.GetFriend(userID, friendID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get friend error, error: %v", err))
		return
	}

	if friend.ID == 0 {
		err = gorm.ErrRecordNotFound
		return
	}
	err = repo.FriendCache.SetFriend(userID, friend)
	return
}

func (c *friendUseCase) CheckFriend(userID string, friendID string) bool {
	friend, err := c.GetFriend(userID, friendID)
	if err != nil {
		return false
	}
	return friend.Status == 1
}

func (c *friendUseCase) CheckBlackFriend(userID string, friendID string) bool {
	friend, err := c.GetFriend(userID, friendID)
	if err == gorm.ErrRecordNotFound {
		return true
	}
	if err != nil {
		return false
	}
	return friend.BlackStatus == 2
}

func (c *friendUseCase) UpdateFriend(userID string, friendID string) (err error) {
	repo.FriendCache.DeleteFriend(userID, friendID)
	return
}

func (c *friendUseCase) GetVersion(userId string) int32 {
	version := repo.FriendCache.GetFriendMaxVersion(userId)
	return int32(version)
}

func (c *friendUseCase) GetUserFriendIdList(userId string) (idList []string) {
	idList = []string{}
	hadFriendIdList, err := db.CloumnList(model.Friend{}, model.Friend{
		OwnerUserID: userId,
		Status:      1,
	}, "friend_user_id")
	if err != nil {
		return
	}
	for _, v := range hadFriendIdList {
		idList = append(idList, v.(string))
	}
	return
}

func (c *friendUseCase) GetFriendRequestInfo(operationID string, userID string, friendID string) (friendRequestInfo *model.FriendRequestInfo, err error) {
	var (
		friendRequest *model.FriendRequest
		user          *userModel.UserBaseInfo
	)

	if friendRequest, err = repo.FriendRequestRepo.Get(userID, friendID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get friend request error, error: %v", err))
		err = code.ErrDB
		return
	}

	if user, err = userUseCase.UserUseCase.GetBaseInfo(userID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get user base info error, error: %v", err))
		err = code.ErrDB
		return
	}

	friendRequestInfo = new(model.FriendRequestInfo)
	util.CopyStructFields(friendRequestInfo, user)
	friendRequestInfo.ID = friendRequest.ID
	friendRequestInfo.ReqMsg = friendRequest.ReqMsg
	friendRequestInfo.CreateTime = friendRequest.CreatedAt
	friendRequestInfo.Status = int64(friendRequest.Status)
	return
}

func (c *friendUseCase) FriendInfoPush(operationID string, userID string, friendID string, mt common.MessageType) (friendInfo *model.FriendInfo, err error) {
	if friendInfo, err = c.GetFriendInfo(userID, friendID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get friend info error, error: %v", err))
		return
	}

	_ = mqtt.SendMessageToUsers(operationID, mt, friendInfo, userID)
	return
}

func (c *friendUseCase) BlackFriendPush(operationID string, userID string, friendID string) {
	blackInfo := model.BlackFriendInfo{
		FromUserId: userID,
		ToUserId:   friendID,
		Msg:        "被好友拒收",
	}
	_ = mqtt.SendMessageToUsers(operationID, common.FriendInBlack, blackInfo, userID)
	return
}

func (c *friendUseCase) FriendRequestInfoPush(operationID string, userID string, friendID string, mt common.MessageType) (friendRequestInfo *model.FriendRequestInfo, err error) {
	if friendRequestInfo, err = c.GetFriendRequestInfo(operationID, userID, friendID); err != nil {
		return
	}

	_ = mqtt.SendMessageToUsers(operationID, mt, friendRequestInfo, friendID)
	return
}

func (c *friendUseCase) CreateFriendLabel(userID, labelID, labelName string) error {

	if isExist, _ := repo.FriendRepo.IsLabelNameExist(userID, labelName); isExist {
		return code.ErrFriendLabelExist
	}

	if _, err := repo.FriendRepo.CreateFriendLabel(userID, labelID, labelName); err != nil {
		return code.ErrDB
	}

	return nil
}

func (c *friendUseCase) DeleteFriendLabel(userID, labelID string) error {

	if isExist, _ := repo.FriendRepo.IsFriendHasLabel(userID, labelID); isExist {
		return code.ErrFriendLabelExist
	}

	if err := repo.FriendRepo.DeleteFriendLabel(userID, labelID); err != nil {
		return code.ErrDB
	}
	return nil
}

func (c *friendUseCase) CheckFriendLabel(userID, labelID string) (isExist bool, err error) {

	if isExist, err = repo.FriendRepo.IsFriendHasLabel(userID, labelID); err != nil {
		return false, code.ErrFriendLabelExist
	}

	return
}

func (c *friendUseCase) UpdateFriendLabel(userID, labelID, labelName string) error {

	if isExist, _ := repo.FriendRepo.IsLabelNameExist(userID, labelName); isExist {
		return code.ErrFriendLabelExist
	}

	if err := repo.FriendRepo.UpdateFriendLabel(userID, labelID, labelName); err != nil {
		return code.ErrDB
	}
	return nil
}

func (c *friendUseCase) GetAllFriendLabels(userID string) ([]*model.FriendLabel, error) {
	friendLabels, err := repo.FriendRepo.GetAllFriendLabels(userID)
	if err != nil {
		return nil, err
	}
	return friendLabels, nil
}

func (c *friendUseCase) ChangeFriendLabel(operateId, userID, labelID string, friendList []string) error {
	batchSize := 50
	total := len(friendList)

	if isExist, _ := repo.FriendRepo.IsFriendLabelExist(userID, labelID); !isExist {
		return code.ErrFriendLabelNotExist
	}

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		batch := friendList[i:end]

		for _, friendID := range batch {

			repo.FriendRepo.ChangeFriendLabel(userID, labelID, friendID)

			c.UpdateFriend(userID, friendID)
		}

		time.Sleep(10 * time.Millisecond)
	}

	mqtt.SendMessageToUsers(operateId, common.FriendInfoLableChage, nil, userID)

	return nil
}

func (c *friendUseCase) GetFriendLabel(userID, Ilabeld string) (*model.FriendLabel, error) {
	friendLabel, err := repo.FriendRepo.GetFriendLabel(userID, Ilabeld)
	if err != nil {
		return nil, err
	}
	return friendLabel, nil
}

func (c *friendUseCase) UpdateFriendBlack(userID, friendID string, status int) error {
	friend, err := c.GetFriend(userID, friendID)
	if err != nil {
		return err
	}
	if friend.BlackStatus == status {
		return nil
	}
	repo.FriendCache.DeleteFriend(userID, friendID)
	return repo.FriendRepo.UpdateBlackStatus(repo.WhereOption{UserId: userID, FriendUserID: friendID}, status)
}
