package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/config"
	chatModel "im/internal/api/chat/model"
	chatUseCase "im/internal/api/chat/usecase"
	"im/internal/api/friend/model"
	"im/internal/api/friend/repo"
	"im/internal/api/friend/usecase"
	groupRepo "im/internal/api/group/repo"
	permissionUseCase "im/internal/api/permission/usecase"
	userModel "im/internal/api/user/model"
	userRepo "im/internal/api/user/repo"
	userUseCase "im/internal/api/user/usecase"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
	"time"
)

var FriendService = new(friendService)

type friendService struct {
}

func (s *friendService) GetLoginUserId(c *gin.Context) (string, error) {
	user_id := c.GetString("user_id")
	return user_id, nil
}

func (s *friendService) Search(c *gin.Context) {
	var (
		req  model.FriendSearchReq
		resp model.FriendSearchResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	req.Check()
	resp.Pagination = req.Pagination
	userID := c.GetString("user_id")

	var users []userModel.User
	if users, resp.Count, err = userRepo.UserRepo.Search(req.Keyword, req.Offset, req.Limit); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("search error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	if len(users) == 0 {
		http.Failed(c, response.GetError(response.ErrFriendSearchNotExist, lang))
		return
	}

	for _, user := range users {
		var friendInfo *model.FriendInfo
		if friendInfo, err = usecase.FriendUseCase.GetFriendInfo(userID, user.UserID); err != nil {
			friendInfo = new(model.FriendInfo)
			util.CopyStructFields(friendInfo, &user)
			friendInfo.UserId = user.UserID
			friendInfo.Status = 2
		}
		resp.List = append(resp.List, *friendInfo)
	}
	http.Success(c, resp)
}

func (s *friendService) AddFriend(c *gin.Context) {
	var (
		req  model.AddFriendReq
		resp model.AddFriendResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if req.UserId == loginUserId {
		http.Failed(c, response.GetError(response.ErrFriendCanNotSelf, lang))
		return
	}

	if usecase.FriendUseCase.CheckFriend(loginUserId, req.UserId) {
		http.Failed(c, response.GetError(response.ErrAlreadyIsFriend, lang))
		return
	}

	var needVerify bool
	if needVerify, err = permissionUseCase.PermissionUseCase.CheckAddFriendPermission(req.OperationID, loginUserId, req.UserId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	if needVerify {

		if _, err = repo.FriendRequestRepo.Create(loginUserId, req.UserId, req.ReqMsg, req.Remark); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("friend request create error, error: %v", err))
			err = response.GetError(response.ErrDB, lang)
			return
		}

		var friendRequestInfo *model.FriendRequestInfo
		if friendRequestInfo, err = usecase.FriendUseCase.FriendRequestInfoPush(req.OperationID, loginUserId, req.UserId, common.FriendRequestPush); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("friend request info push error, error: %v", err))
			err = response.GetError(response.ErrUnknown, lang)
			return
		}

		util.CopyStructFields(&resp, friendRequestInfo)
	} else {
		if err = usecase.FriendUseCase.AddFriend(req.OperationID, loginUserId, req.UserId, req.Remark, "", false); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("add friend error, error: %v", err))
			err = response.GetError(response.ErrUnknown, lang)
			return
		}
	}

	http.Success(c, resp)
}

func (s *friendService) AddFriendAck(c *gin.Context) {
	var (
		req  model.AddFriendAckReq
		resp model.AddFriendAckResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	var friendRequest *model.FriendRequest
	if friendRequest, err = repo.FriendRequestRepo.GetByID(req.ReqId); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get by id error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrApplyNotFound, lang))
		return
	}

	if friendRequest.Status != 0 {
		http.Failed(c, response.GetError(response.ErrApplyDone, lang))
		return
	}

	if friendRequest.ToUserID != loginUserId {
		http.Failed(c, response.GetError(response.ErrNoPermission, lang))
		return
	}

	if req.Status == 1 {

		if err = usecase.FriendUseCase.AddFriend(req.OperationID, friendRequest.FromUserID, friendRequest.ToUserID, friendRequest.Remark, "", false); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("add friend error, error: %v", err))
			http.Failed(c, response.GetError(response.ErrDB, lang))
			return
		}

		if err = repo.FriendRequestRepo.UpdateStatus(friendRequest.FromUserID, friendRequest.ToUserID, 1); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update status error, error: %v", err))
			http.Failed(c, response.GetError(response.ErrDB, lang))
			return
		}
	} else {

		if err = repo.FriendRequestRepo.UpdateStatus(friendRequest.FromUserID, friendRequest.ToUserID, 2); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("update status error, error: %v", err))
			http.Failed(c, response.GetError(response.ErrDB, lang))
			return
		}
	}

	var friendRequestInfo *model.FriendRequestInfo
	if friendRequestInfo, err = usecase.FriendUseCase.FriendRequestInfoPush(req.OperationID, friendRequest.FromUserID, friendRequest.ToUserID, common.AddFriendAckPush); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("friend request info push error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	util.CopyStructFields(&resp, friendRequestInfo)
	http.Success(c, resp)
}

func (s *friendService) DeleteFriend(c *gin.Context) {
	var (
		req  model.DeleteFriendReq
		resp model.DeleteFriendResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckDeleteFriendPermission(req.OperationID, loginUserId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	hadFriend := model.Friend{}

	err = db.Info(&hadFriend, model.Friend{
		OwnerUserID:  loginUserId,
		FriendUserID: req.UserId,
		Status:       1,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	hadToFriend := model.Friend{}

	err = db.Info(&hadToFriend, model.Friend{
		OwnerUserID:  req.UserId,
		FriendUserID: loginUserId,
		Status:       1,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	err = usecase.FriendUseCase.DeleteFriend(req.OperationID, loginUserId, req.UserId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	db.Info(&hadFriend, hadFriend.ID)
	res, _ := usecase.FriendUseCase.GetFriendInfo(loginUserId, req.UserId)
	util.CopyStructFields(&resp, res)
	http.Success(c, resp)
}

func (s *friendService) SetFriendRemark(c *gin.Context) {
	var (
		req  model.SetFriendRemarkReq
		resp model.SetFriendRemarkResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	hadFriend := model.Friend{}

	err = db.Info(&hadFriend, model.Friend{
		OwnerUserID:  loginUserId,
		FriendUserID: req.UserId,
		Status:       1,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if err = repo.FriendRepo.UpdateFriendRemark(repo.WhereOption{UserId: loginUserId, FriendUserID: req.UserId}, &model.Friend{Remark: req.Remark}); err != nil {
		logger.Sugar.Errorf("OperationID:%s,err:%s", req.OperationID, err.Error())
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}

	usecase.FriendUseCase.UpdateFriend(loginUserId, req.UserId)

	friendInfo, _ := usecase.FriendUseCase.FriendInfoPush(req.OperationID, loginUserId, req.UserId, common.ChangeFriendPush)

	util.CopyStructFields(&resp, friendInfo)
	http.Success(c, resp)
}

func (s *friendService) GetFriendRemark(c *gin.Context) {
	var (
		req     model.GetFriendRemarkReq
		resp    model.GetFriendRemarkResp
		friends []model.Friend
		err     error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if friends, err = repo.FriendRepo.GetFriendRemark(loginUserId); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	for _, friend := range friends {
		resp.Remark = append(resp.Remark, friend.Remark)
	}

	http.Success(c, resp)
}

func (s *friendService) CheckFriendRemark(c *gin.Context) {
	var (
		req     model.CheckFriendRemarkReq
		err     error
		isExist bool
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	hadFriend := model.Friend{}

	err = db.Info(&hadFriend, model.Friend{
		OwnerUserID:  loginUserId,
		FriendUserID: req.UserId,
		Status:       1,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if isExist, err = repo.FriendRepo.IsFriendHasRemark(loginUserId, req.UserId, req.Remark); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	if !isExist {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}

	http.Success(c)
}

func (s *friendService) GetFriendsInfo(c *gin.Context) {
	var (
		req  model.GetFriendInfoReq
		resp model.GetFriendInfoResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	friend := model.Friend{}
	err = db.Info(&friend, model.Friend{
		OwnerUserID:  loginUserId,
		FriendUserID: req.UserId,
		Status:       1,
	})
	if err != nil {

	}
	friendInfo, err := usecase.FriendUseCase.GetFriendInfo(loginUserId, req.UserId)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " GetFriendsInfo GetFriendMemberUserInfo error:", err.Error())
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	currentUser, err1 := userUseCase.UserUseCase.GetInfo(c.GetString("user_id"))
	if err1 != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " GetFriendsInfo UserUseCase GetInfo error:", err.Error())
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}

	if currentUser.IsPrivilege != 1 {
		friendInfo.PhoneNumber = ""
	}
	if err = util.CopyStructFields(&resp, friendInfo); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "GetFriendsInfo CopyStructFields error:", err.Error())
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	username := fmt.Sprintf("%s_%s", config.Config.Station, req.UserId)
	onlineClients, err := mqtt.GetClients(username, "", "", mqtt.ConnStateTypeConnected, 1, 1)
	if err != nil {
		user, err := userUseCase.UserUseCase.GetInfo(req.UserId)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user GetInfo, error: %v", err))
			http.Failed(c, response.GetError(response.ErrUserIdNotExist, lang))
			return
		}
		resp.Online = false
		timeLayout := "2006/01/02 15:04:05"
		timeStr := time.Unix(user.LatestLoginTime, 0).Format(timeLayout)
		resp.OfflineInfo = fmt.Sprintf("离线 %s", timeStr)
	} else {
		user, _ := userUseCase.UserUseCase.GetInfo(req.UserId)
		for i := len(onlineClients) - 1; i >= 0; i-- {
			v := onlineClients[i]
			temp := strings.Split(v.ClientID, "_")
			if len(temp) < 2 {
				continue
			}
			resp.Online = true

			resp.LoginIp = user.LoginIp
			resp.LoginIpLocaltion, _ = util.QueryIpRegion(resp.LoginIp)
			break
		}
		if !resp.Online {
			user, err = userUseCase.UserUseCase.GetInfo(req.UserId)
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("user GetInfo, error: %v", err))
				http.Failed(c, response.GetError(response.ErrUserIdNotExist, lang))
				return
			}
			resp.Online = false
			timeLayout := "2006/01/02 15:04:05"
			timeStr := time.Unix(user.LatestLoginTime, 0).Format(timeLayout)
			resp.OfflineInfo = fmt.Sprintf("离线 %s", timeStr)
		}
	}

	http.Success(c, resp)
}

func (s *friendService) GetFriendsMsgMaxSeq(c *gin.Context) {
	var (
		req  model.GetFriendListMaxSeqReq
		resp model.GetFriendListMaxSeqResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 20
	}
	total := int64(0)
	friends := []model.Friend{}
	wheres := map[string]interface{}{
		"owner_user_id": loginUserId,
		"status":        1,
	}

	if err = db.Find(model.Friend{}, wheres, "", req.Page, req.PageSize, &total, &friends); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get friend error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}

	for _, friend := range friends {
		friendsSeq := model.MaxFriendsSeq{}
		conversationID := chatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, loginUserId, friend.FriendUserID)
		friendsSeq.ConversationID = conversationID
		resp.List = append(resp.List, friendsSeq)
	}

	resp.Count = total
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(c, resp)
}

func (s *friendService) GetFriendList(c *gin.Context) {
	var (
		req  model.GetFriendListReq
		resp model.GetFriendListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 20
	}
	total := int64(0)
	friends := []model.Friend{}
	wheres := map[string]interface{}{
		"owner_user_id": loginUserId,
		"status":        1,
	}
	if req.SearchKey != "" {
		uids := groupRepo.GroupRepo.GetUserIdListByNickName(req.SearchKey)
		wheres["friend_user_id"] = uids
	}
	if req.FriendLabel != "" {
		wheres["friend_label"] = req.FriendLabel
	}
	if req.BlackStatus == 1 {
		wheres["black_status"] = req.BlackStatus
	}
	if req.Version > 0 {
		wheres["version"] = fmt.Sprintf(">=%d", req.Version)
	}

	if err = db.Find(model.Friend{}, wheres, "", req.Page, req.PageSize, &total, &friends); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("get friend error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}

	for _, friend := range friends {
		friendInfo, _ := usecase.FriendUseCase.GetFriendInfo(loginUserId, friend.FriendUserID)
		resp.List = append(resp.List, *friendInfo)
	}

	resp.Count = total
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *friendService) GetBlackList(c *gin.Context) {
	var (
		req  model.GetBlackListReq
		resp model.GetBlackListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	req.Pagination.Check()
	if resp.List, resp.Count, err = repo.FriendRepo.FetchBlackUserList(loginUserId, req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *friendService) GetBlackListV2(c *gin.Context) {
	var (
		req  model.GetBlackListReq
		resp model.GetBlackListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	req.Pagination.Check()
	wheres := map[string]interface{}{
		"owner_user_id": loginUserId,
		"status":        1,
		"black_status":  model.InBlack,
	}
	var friends []model.Friend
	if err = db.Find(model.Friend{}, wheres, "", req.Page, req.PageSize, &req.Count, &friends); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("error: %v", err))
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}

	for _, friend := range friends {
		user, err1 := userUseCase.UserUseCase.GetInfo(friend.FriendUserID)
		if err1 != nil {
			continue
		}
		fD := model.FriendInfo{}
		_ = util.Copy(user, &fD)
		resp.List = append(resp.List, fD)
	}

	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *friendService) GetFriendApplyList(c *gin.Context) {
	var (
		req    model.GetFriendApplyListReq
		resp   model.GetFriendApplyListRes
		err    error
		applys []model.FriendRequest
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)
	total := int64(0)
	if err = db.Find(model.FriendRequest{}, map[string]interface{}{"to_user_id": loginUserId}, "id desc", req.Page, req.PageSize, &total, &applys); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), " GetFriendApplyList find error:", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if resp.List, err = repo.FriendRequestRepo.GetFriendMemberUserInfoList(&applys, "from_user_id"); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), " GetFriendApplyList GetFriendMemberUserInfoList error:", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	resp.Count = total
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(c, resp)
}

func (s *friendService) GetFriendApplyAll(c *gin.Context) {
	var (
		req    model.GetFriendApplyListReq
		resp   model.GetFriendApplyListRes
		err    error
		applys []model.FriendRequest
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)
	total := int64(0)
	if err = db.Find(model.FriendRequest{}, map[string]interface{}{"to_user_id|from_user_id": loginUserId}, "id desc", req.Page, req.PageSize, &total, &applys); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), " GetFriendApplyList find error:", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if resp.List, err = repo.FriendRequestRepo.GetFriendMemberUserInfoAll(&applys, loginUserId); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), " GetFriendApplyList GetFriendMemberUserInfoList error:", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	resp.Count = total
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(c, resp)
}

func (s *friendService) GetSelfFriendApplyList(c *gin.Context) {
	var (
		req    model.GetFriendApplyListReq
		resp   model.GetFriendApplyListRes
		err    error
		applys []model.FriendRequest
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)
	total := int64(0)
	if err = db.Find(model.FriendRequest{}, map[string]interface{}{"from_user_id": loginUserId}, "id desc", req.Page, req.PageSize, &total, &applys); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), " GetSelfFriendApplyList find error:", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if resp.List, err = repo.FriendRequestRepo.GetFriendMemberUserInfoList(&applys, "to_user_id"); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), " GetSelfFriendApplyList GetFriendMemberUserInfoList error:", err)
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	resp.Count = total
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	http.Success(c, resp)
}

func (s *friendService) GetFriendMaxSeq(c *gin.Context) {
	var (
		req  model.GetFriendMaxReq
		resp model.MaxSeq
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	maxSeq, maxUpdateTime := repo.FriendRepo.FriendMaxSeq(repo.WhereOption{UserId: loginUserId})
	resp.Version = maxSeq
	resp.MaxUpdateTime = maxUpdateTime
	http.Success(c, resp)
}

func (s *friendService) IsFriend(c *gin.Context) {
	var (
		req    model.IsFriendReq
		friend *model.Friend
		err    error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	friend, err = repo.FriendRepo.GetFriend(loginUserId, req.FriendUserID)
	if err != nil || friend.ID == 0 {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), fmt.Sprintf("error: %v", err))
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	http.Success(c)
}

func (s *friendService) SearchFriend(c *gin.Context) {
	var (
		req        model.SearchFriendReq
		resp       model.SearchFriendRes
		friendInfo *model.FriendInfo
		err        error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	hadUser, err := userUseCase.UserUseCase.GetBaseInfo(req.UserId)
	if err != nil {
		logger.Sugar.Debug(req.OperationID, fmt.Sprintf("查询id失败 err: %v", err))
		hadUser, err = userUseCase.UserUseCase.GetBaseInfoByPhoneNumber(req.UserId)
		logger.Sugar.Debug(req.OperationID, fmt.Sprintf("查询手机号结果 %v err: %v", hadUser, err))
	}
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFriendSearchNotExist, lang))
		return
	}
	friend := model.Friend{}
	err = db.Info(&friend, model.Friend{
		OwnerUserID:  loginUserId,
		FriendUserID: hadUser.UserId,
		Status:       1,
	})

	if err != nil {
		friendInfo = new(model.FriendInfo)
		util.CopyStructFields(friendInfo, hadUser)
		friendInfo.UserId = hadUser.UserId
		friendInfo.Status = 2
	} else {

		if friendInfo, err = usecase.FriendUseCase.GetFriendInfo(loginUserId, req.UserId); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("db get friend error, error: %v", err))
			http.Failed(c, response.GetError(response.ErrFriendSearchNotExist, lang))
			return
		}
	}
	currentUser, err1 := userUseCase.UserUseCase.GetInfo(c.GetString("user_id"))
	if err1 != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), " UserUseCase GetInfo error:", err.Error())
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}

	if currentUser.IsPrivilege != 1 {
		friendInfo.PhoneNumber = ""
	}
	if err = util.CopyStructFields(&resp, friendInfo); err != nil {
		logger.Sugar.Error(req.OperationID, "func", util.GetSelfFuncName(), "GetFriendsInfo CopyStructFields error:", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	http.Success(c, resp)
}

func (s *friendService) FriendListSync(c *gin.Context) {
	var (
		req  model.FriendListSyncReq
		resp model.FriendListSyncResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 20
	}

	total := int64(0)
	friends := []model.Friend{}
	wheres := map[string]interface{}{
		"owner_user_id": loginUserId,
		"version":       fmt.Sprintf(">%d", req.Version),
	}

	db.Find(model.Friend{}, wheres, "version desc", req.Page, req.PageSize, &total, &friends)

	for _, friend := range friends {
		friendInfo, _ := usecase.FriendUseCase.GetFriendInfo(loginUserId, friend.FriendUserID)
		resp.List = append(resp.List, *friendInfo)
	}

	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Count = total
	http.Success(c, resp)
}

func (s *friendService) CreateFriendLabel(c *gin.Context) {
	var (
		req  model.CreateFriendLabelReq
		resp model.CreateFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	labelId := util.RandID(10)
	err = usecase.FriendUseCase.CreateFriendLabel(loginUserId, labelId, req.LabelName)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("创建分组失败 error: %v", err))
		http.Failed(c, err)
		return
	}
	resp.LabelId = labelId
	resp.LabelName = req.LabelName
	http.Success(c, resp)
}

func (s *friendService) DeleteFriendLabel(c *gin.Context) {
	var (
		req  model.DeleteFriendLabelReq
		resp model.DeleteFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if req.LabelId == loginUserId {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("默认好友分组不可修改 error: %v", err))
		http.Failed(c, response.GetError(response.ErrFriendLabelDelete, lang))
		return
	}
	err = usecase.FriendUseCase.DeleteFriendLabel(loginUserId, req.LabelId)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("删除分组失败 error: %v", err))
		http.Failed(c, err)
		return
	}

	http.Success(c, resp)
}

func (s *friendService) UpdateFriendLabel(c *gin.Context) {
	var (
		req  model.UpdateFriendLabelReq
		resp model.UpdateFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if req.LabelId == loginUserId {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("默认好友分组不可修改 error: %v", err))
		http.Failed(c, response.GetError(response.ErrFriendLabelForbiden, lang))
		return
	}
	err = usecase.FriendUseCase.UpdateFriendLabel(loginUserId, req.LabelId, req.LabelName)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("删除分组失败 error: %v", err))
		http.Failed(c, err)
		return
	}
	resp.LabelId = req.LabelId
	resp.LabelName = req.LabelName
	http.Success(c, resp)
}

func (s *friendService) GetFriendLabel(c *gin.Context) {
	var (
		req  model.GetFriendLabelReq
		resp model.GetFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	friendLabels, err := usecase.FriendUseCase.GetAllFriendLabels(loginUserId)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("获取好友分组信息失败 error: %v", err))
		http.Failed(c, err)
		return
	}

	for _, v := range friendLabels {
		temp := model.FriendLabelInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}

	http.Success(c, resp)
}

func (s *friendService) ChangeFriendLabel(c *gin.Context) {
	var (
		req  model.ChangeFriendLabelReq
		resp model.ChangeFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if err = usecase.FriendUseCase.ChangeFriendLabel(req.OperationID, loginUserId, req.LabelId, req.FriendList); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("修改好友分组 error: %v", err))
		http.Failed(c, err)
		return
	}

	http.Success(c, resp)
}

func (s *friendService) AddBlack(c *gin.Context) {
	var (
		req  model.ParamsCommFriend
		resp model.ChangeFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if err = usecase.FriendUseCase.UpdateFriendBlack(loginUserId, req.UserId, model.InBlack); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("添加黑名 error: %v", err))
		http.Failed(c, err)
		return
	}

	http.Success(c, resp)
}

func (s *friendService) RemoveBlack(c *gin.Context) {
	var (
		req  model.ParamsCommFriend
		resp model.ChangeFriendLabelResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("bind json error: %v", err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	if err = usecase.FriendUseCase.UpdateFriendBlack(loginUserId, req.UserId, model.NotBlack); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), fmt.Sprintf("移除黑名单 error: %v", err))
		http.Failed(c, err)
		return
	}

	http.Success(c, resp)
}
