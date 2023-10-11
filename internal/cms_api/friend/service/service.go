package service

import (
	"fmt"
	chatModel "im/internal/api/chat/model"
	chatRepo "im/internal/api/chat/repo"
	chatUseCase "im/internal/api/chat/usecase"
	"im/internal/api/friend/model"
	apiModel "im/internal/api/friend/model"
	apiFriendUseCase "im/internal/api/friend/usecase"
	apiUserModel "im/internal/api/user/model"
	apiUserRepo "im/internal/api/user/repo"
	cmsModel "im/internal/cms_api/friend/model"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"

	"github.com/gin-gonic/gin"
)

var FriendService = new(friendService)

type friendService struct{}

func (s *friendService) GetLoginUserId(c *gin.Context) (string, error) {

	userId := c.GetString("o_user_id")
	return userId, nil
}

func (s *friendService) UserFriendList(c *gin.Context) {
	var (
		req  cmsModel.UserFriendListReq
		resp cmsModel.UserFriendListResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	friendList := []apiModel.Friend{}
	total := int64(0)
	db.Find(apiModel.Friend{}, apiModel.Friend{
		OwnerUserID:  req.UserId,
		FriendUserID: req.FriendId,
		Status:       1,
	}, "id desc", req.Page, req.PageSize, &total, &friendList)

	for _, friend := range friendList {
		friendInfo, _ := apiFriendUseCase.FriendUseCase.GetFriendInfo(req.UserId, friend.FriendUserID)
		resp.List = append(resp.List, *friendInfo)
	}

	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Count = int(total)
	http.Success(c, resp)
}

func (s *friendService) UserAddFriend(c *gin.Context) {
	var (
		req  cmsModel.UserAddFriendReq
		resp cmsModel.UserAddFriendResp
		user *apiUserModel.User
		err  error
	)
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	logger.Sugar.Infof("OperationID:%s,data:%v", req.OperationID, req)

	opt := apiUserRepo.WhereOption{
		UserId:      req.FriendId,
		NickName:    req.FriendId,
		PhoneNumber: req.FriendId,
	}
	if user, err = apiUserRepo.UserRepo.GetByFriendID(opt); err != nil {
		http.Failed(c, code.ErrUserNotFound)
		return
	}
	if user.Status == constant.UserStatusFreeze {
		http.Failed(c, code.ErrUserFreeze)
		return
	}
	if user.Status != constant.UserStatusNormal {
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	req.FriendId = user.UserID
	hadFriend := apiModel.Friend{}
	err = db.Info(&hadFriend, apiModel.Friend{
		OwnerUserID:  req.UserId,
		FriendUserID: req.FriendId,
		Status:       1,
	})
	if err == nil {
		http.Success(c, resp)
		return
	}

	tx := db.DB.Begin()

	err = db.DeleteTx(tx, model.Friend{}, model.Friend{
		OwnerUserID:  req.UserId,
		FriendUserID: req.FriendId,
	})
	if err != nil {
		logger.Sugar.Errorf("OperationID:%s,err:%s", req.OperationID, err.Error())
		http.Failed(c, code.ErrFailRequest)
		return
	}

	err = db.DeleteTx(tx, model.Friend{}, model.Friend{
		OwnerUserID:  req.FriendId,
		FriendUserID: req.UserId,
	})
	if err != nil {
		logger.Sugar.Errorf("OperationID:%s,err:%s", req.OperationID, err.Error())
		http.Failed(c, code.ErrFailRequest)
		return
	}

	data := model.Friend{
		OwnerUserID:    req.UserId,
		FriendLabel:    req.UserId,
		FriendUserID:   req.FriendId,
		OperatorUserID: "",
		Status:         1,
	}
	if err = db.InsertTx(tx, &data); err != nil {
		logger.Sugar.Errorf("OperationID:%s,err:%s", req.OperationID, err.Error())
		http.Failed(c, code.ErrFailRequest)
		return
	}
	toFriendData := model.Friend{
		OwnerUserID:    req.FriendId,
		FriendLabel:    req.FriendId,
		FriendUserID:   req.UserId,
		OperatorUserID: "",
		Status:         1,
	}
	if err = db.InsertTx(tx, &toFriendData); err != nil {
		logger.Sugar.Errorf("OperationID:%s,err:%s", req.OperationID, err.Error())
		http.Failed(c, code.ErrFailRequest)
		return
	}
	tx.Commit()

	apiFriendUseCase.FriendUseCase.UpdateFriend(req.UserId, req.FriendId)
	apiFriendUseCase.FriendUseCase.UpdateFriend(req.FriendId, req.UserId)

	apiFriendUseCase.FriendUseCase.FriendInfoPush(req.OperationID, req.UserId, req.FriendId, common.AddFriendPush)
	apiFriendUseCase.FriendUseCase.FriendInfoPush(req.OperationID, req.FriendId, req.UserId, common.AddFriendPush)

	var (
		message *chatModel.MessageInfo
		content chatModel.MessageContent
	)
	content.OperatorID = req.UserId

	conversationID := chatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, req.UserId, req.FriendId)
	if message, err = chatUseCase.MessageUseCase.SendSystemMessageToUsers(req.OperationID, chatModel.ConversationTypeSingle, conversationID, chatModel.MessageFriendAddNotify, &content, req.UserId, req.FriendId); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("send system message to users error, error: %v", err))
	}

	if err = chatUseCase.ConversationUseCase.UpsertUsersStartSeq(req.OperationID, chatModel.ConversationTypeSingle, conversationID, message.Seq-1, req.UserId, req.FriendId); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("upsert users start seq error, error: %v", err))
	}

	http.Success(c, resp)
	return
}

func (s *friendService) DeleteFriend(c *gin.Context) {
	var (
		req  cmsModel.UserRemoveFriendReq
		resp cmsModel.UserRemoveFriendResp
		err  error
	)
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	hadFriend := model.Friend{}

	err = db.Info(&hadFriend, model.Friend{
		OwnerUserID:  req.UserId,
		FriendUserID: req.FriendId,
		Status:       1,
	})
	if err != nil {
		http.Failed(c, code.ErrFriendNotExist)
		return
	}
	hadToFriend := model.Friend{}

	err = db.Info(&hadToFriend, model.Friend{
		OwnerUserID:  req.FriendId,
		FriendUserID: req.UserId,
		Status:       1,
	})
	if err != nil {
		http.Failed(c, code.ErrFriendNotExist)
		return
	}
	tx := db.DB.Begin()
	if err = db.UpdateTx(tx, model.Friend{}, model.Friend{
		OwnerUserID:  req.UserId,
		FriendUserID: req.FriendId,
	}, model.Friend{Status: 2}); err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}

	if err = db.UpdateTx(tx, model.Friend{}, model.Friend{
		OwnerUserID:  req.FriendId,
		FriendUserID: req.UserId,
	}, model.Friend{Status: 2}); err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	tx.Commit()

	apiFriendUseCase.FriendUseCase.UpdateFriend(hadFriend.OwnerUserID, hadFriend.FriendUserID)
	apiFriendUseCase.FriendUseCase.UpdateFriend(hadFriend.FriendUserID, hadFriend.OwnerUserID)

	conversationID := chatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, req.UserId, req.FriendId)
	if err2 := chatRepo.ConversationRepo.Delete(chatModel.ConversationTypeSingle, conversationID); err2 != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("delete single conversation error, conversation id: %s, error: %v", conversationID, err))
	}

	apiFriendUseCase.FriendUseCase.FriendInfoPush(req.OperationID, hadFriend.OwnerUserID, hadFriend.FriendUserID, common.RemoveFriendPush)
	apiFriendUseCase.FriendUseCase.FriendInfoPush(req.OperationID, hadFriend.FriendUserID, hadFriend.OwnerUserID, common.RemoveFriendPush)
	http.Success(c, resp)
}
