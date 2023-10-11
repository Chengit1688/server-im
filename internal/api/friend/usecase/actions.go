package usecase

import (
	"encoding/json"
	"fmt"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	chatUseCase "im/internal/api/chat/usecase"
	"im/internal/api/friend/model"
	"im/internal/api/friend/repo"
	momentsModel "im/internal/api/moments/model"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/util"
	"time"
)

func (c *friendUseCase) UserInfoChange(userID string) error {
	var err error
	friends := []model.Friend{}
	wheres := map[string]interface{}{
		"owner_user_id": userID,
		"status":        1,
	}
	total := int64(0)
	if err = db.Find(model.Friend{}, wheres, "", 0, 1000, &total, &friends); err != nil {
		logger.Sugar.Errorw("func", util.GetSelfFuncName(), "error", fmt.Sprintf("get friend error, error: %v", err))
		return err
	}
	total_ := int(total)
	batchSize := 50
	for i := 0; i < total_; i += batchSize {
		end := i + batchSize
		if end > total_ {
			end = total_
		}

		batch := friends[i:end]

		for _, friendID := range batch {

			repo.FriendRepo.ChangeFriendVersion(userID, friendID.FriendUserID)

			c.UpdateFriend(userID, friendID.FriendUserID)
		}

		time.Sleep(10 * time.Millisecond)
	}
	return err
}

func (c *friendUseCase) AddFriend(operationID string, userID string, friendID string, remark string, greeting string, isCustomer bool) (err error) {
	if userID == "" || friendID == "" {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user id: %s, friend id: %s", userID, friendID))
		err = code.ErrUserNotFound
		return
	}

	if userID == friendID {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user id = friend id"))
		err = code.ErrFriendCanNotSelf
		return
	}

	if err = repo.FriendRepo.AddFriend(userID, friendID, userID, remark, 0); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user add friend error, error: %v", err))
		err = code.ErrDB
		return
	}

	if err = repo.FriendRepo.AddFriend(friendID, userID, userID, "", 0); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("friend add user error, error: %v", err))
		err = code.ErrDB
		return
	}

	if !isCustomer {
		if err = repo.FriendRequestRepo.UpdateStatus(userID, friendID, 1); err != nil {
			logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update status error, error: %v", err))
			err = code.ErrDB
			return
		}
	}

	c.UpdateFriend(userID, friendID)
	c.UpdateFriend(friendID, userID)

	c.FriendInfoPush(operationID, userID, friendID, common.AddFriendPush)
	c.FriendInfoPush(operationID, friendID, userID, common.AddFriendPush)

	var (
		message *chatModel.MessageInfo
		content chatModel.MessageContent
	)
	content.OperatorID = userID
	conversationID := chatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, userID, friendID)
	if message, err = chatUseCase.MessageUseCase.SendSystemMessageToUsers(operationID, chatModel.ConversationTypeSingle, conversationID, chatModel.MessageFriendAddNotify, &content, userID, friendID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send system message to users error, error: %v", err))
		return err
	}

	if err = chatUseCase.ConversationUseCase.UpsertUsersStartSeq(operationID, chatModel.ConversationTypeSingle, conversationID, message.Seq-1, userID, friendID); err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("upsert users start seq error, error: %v", err))
		return err
	}

	if greeting != "" {
		greetingContent := chatModel.MsgContent{}
		greetingContent.Text = greeting
		contentData, _ := json.Marshal(&greetingContent)
		contentEncrypt, _ := util.Encrypt(contentData, common.ContentKey)
		chatUseCase.MessageUseCase.SendMessageToUsers(operationID, "", chatModel.ConversationTypeSingle, conversationID, friendID, chatModel.MessageText, contentEncrypt, friendID, userID)
	}
	return err
}

func (c *friendUseCase) DeleteFriend(operationID string, FromUserId, ToUserId string) (err error) {
	tx := db.DB.Begin()
	if err = db.UpdateTx(tx, model.Friend{}, model.Friend{
		OwnerUserID:  FromUserId,
		FriendUserID: ToUserId,
	}, model.Friend{Status: 2}); err != nil {
		return
	}

	if err = db.UpdateTx(tx, model.Friend{}, model.Friend{
		OwnerUserID:  ToUserId,
		FriendUserID: FromUserId,
	}, model.Friend{Status: 2}); err != nil {
		return
	}

	if err = tx.Where("publisher_user_id = ?", ToUserId).
		Where("friend_user_id = ?", FromUserId).
		Delete(&momentsModel.MomentsInbox{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	_ = c.UpdateFriend(FromUserId, ToUserId)
	_ = c.UpdateFriend(ToUserId, FromUserId)

	conversationID := chatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, FromUserId, ToUserId)
	if err2 := chatRepo.ConversationRepo.Delete(chatModel.ConversationTypeSingle, conversationID); err2 != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("delete users conversation error, user id: %s, friend user id: %s, error: %v", FromUserId, ToUserId, err))
	}
	hadFriend := model.Friend{}

	_ = db.Info(&hadFriend, model.Friend{
		OwnerUserID:  FromUserId,
		FriendUserID: ToUserId,
	})

	hadToFriend := model.Friend{}

	_ = db.Info(&hadToFriend, model.Friend{
		OwnerUserID:  ToUserId,
		FriendUserID: FromUserId,
	})

	c.FriendInfoPush(operationID, FromUserId, ToUserId, common.RemoveFriendPush)
	c.FriendInfoPush(operationID, ToUserId, FromUserId, common.RemoveFriendPush)
	return nil
}
