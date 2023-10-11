package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	chatModel "im/internal/api/chat/model"
	apiChatUseCase "im/internal/api/chat/usecase"
	friendusecase "im/internal/api/friend/usecase"
	"im/internal/api/moments/model"
	momentsRepo "im/internal/api/moments/repo"
	"im/internal/api/moments/usecase"
	userRepo "im/internal/api/user/repo"
	userusecase "im/internal/api/user/usecase"
	"im/pkg/common"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
	"time"
)

var MomentsService = new(momentsService)

type momentsService struct{}

func (s *momentsService) AddMomentsMessage(c *gin.Context) {
	var (
		req  model.IssueReq
		resp model.MomentsMessage
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	logger.Sugar.Infof("OperationID:%s,data:%v", req.OperationID, req)
	userID := c.GetString("user_id")

	add := model.MomentsMessage{
		Content:        req.Content,
		ShareTagID:     util.Int64SliceToString(req.ShareTagID, ","),
		ShareFriendID:  strings.Join(req.ShareFriendID, ","),
		InviteFriendID: strings.Join(req.InviteFriendID, ","),
		Location:       req.Location,
		NoComment:      req.NoComment,
		CanSee:         req.CanSee,
		Image:          strings.Join(req.Images, ";"),
		Video:          strings.Join(req.Videos, ";"),
		VideoImg:       strings.Join(req.VideoImg, ";"),
		UserId:         userID,
		Year:           time.Now().Year(),
		Month:          util.StringToInt(time.Now().Format("01")),
		CreatedAt:      time.Now().Unix(),
		Day:            time.Now().Day(),
	}
	if req.CanSee == 0 {
		add.CanSee = model.CanSeePublic
	}
	if add, err = momentsRepo.MomentsRepo.MomentsAdd(add); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	for _, inviter := range req.InviteFriendID {
		userInfo, err1 := userusecase.UserUseCase.GetBaseInfo(userID)
		if err1 != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			continue
		}
		conversationID := apiChatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, userID, inviter)
		pushData := model.MomentsInviteFriendPush{
			ConversationID:    conversationID,
			Timestamp:         time.Now().Unix(),
			PublisherUserID:   userID,
			FriendUserID:      inviter,
			MomentsID:         add.ID,
			FaceURL:           userInfo.FaceURL,
			PublisherNickname: userInfo.NickName,
		}
		if err = mqtt.SendMessageToUsers(req.OperationID, common.MomentsInviteFriend, pushData, []string{userID, inviter}...); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		}
	}
	_ = util.Copy(add, &resp)
	userRepo.UserCache.DelMomentsMessage(userID)

	http.Success(c, add)
}

func (s *momentsService) GetMomentsMessage(c *gin.Context) {
	var (
		req      model.IssueListReq
		resp     model.IssueList
		err      error
		count    int64
		moments  []model.MomentsMessage
		info     model.IssueInfo
		list     []model.IssueInfo
		strCount string
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")

	req.IsOwner = 1
	userKey := req.UserID
	if req.UserID == "" {
		req.IsOwner = 2
		userKey = userID
	}
	req.Check()
	resp.Pagination = req.Pagination

	list, _, strCount = userRepo.UserCache.GetMomentsMessage(userKey, req.IsOwner, req.Page)
	if strCount != "" {
		count = util.StringToInt64(strCount)
	}

	if len(list) == 0 {
		if req.UserID != "" {

			moments, err = momentsRepo.MomentsRepo.GetSelfMomentsList(userID, req, &count)
		} else {

			req.UserID = userID
			moments, err = momentsRepo.MomentsRepo.GetFriendMomentsList(req, &count)
		}
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
				http.Failed(c, response.GetError(response.ErrFailRequest, lang))
				return
			}
		}
		if len(moments) > 0 {
			for _, v := range moments {
				if info, err = usecase.MomentsUseCase.MomentsDetail(v, lang, resp.Pagination.Offset, resp.Pagination.Limit); err != nil {
					if err != gorm.ErrRecordNotFound {
						logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
						http.Failed(c, err)
						return
					}
				}
				list = append(list, info)
			}
		}
		if err = userRepo.UserCache.SetMomentsMessage(userKey, req.IsOwner, req.Page, list, count); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
	}
	resp.Count = count
	resp.List = list
	http.Success(c, resp)
}

func (s *momentsService) DelMomentsMessage(c *gin.Context) {
	var (
		req  model.DelIssueReq
		resp model.MomentsMessage
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userId := c.GetString("user_id")
	user, err := userusecase.UserUseCase.GetBaseInfo(userId)
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	if resp, err = momentsRepo.MomentsRepo.MomentsDel(req, userId, user.IsPrivilege); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserPermissions, lang))
		return
	}
	userRepo.UserCache.DelMomentsMessage(userId)
	http.Success(c, resp)
	return
}

func (s *momentsService) MomentsCommentsAdd(c *gin.Context) {
	var (
		req  model.CommentReq
		resp model.MomentsComments
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	info := new(model.MomentsMessage)
	if err = momentsRepo.MomentsRepo.MomentsInfo(req.MomentsID, info); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	add := new(model.MomentsComments)
	util.CopyStructFields(add, &req)
	if req.Images != nil && len(req.Images) > 0 {
		add.Image = strings.Join(req.Images, ";")
	}
	userID := c.GetString("user_id")
	if userID == info.UserId {
		add.IsOwnComment = 2
	}
	add.UserID = userID
	if err = momentsRepo.MomentsRepo.MomentsCommentsAdd(*add); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	publisherUser, err := userusecase.UserUseCase.GetBaseInfo(add.UserID)
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	resp = *add
	userRepo.UserCache.DelMomentsMessage(info.UserId)
	conversationID := apiChatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, add.UserID, info.UserId)
	pushData := model.MomentsInviteFriendPush{
		ConversationID:    conversationID,
		Timestamp:         time.Now().Unix(),
		PublisherUserID:   add.UserID,
		FriendUserID:      info.UserId,
		MomentsID:         req.MomentsID,
		FaceURL:           publisherUser.FaceURL,
		PublisherNickname: publisherUser.NickName,
	}
	if err = mqtt.SendMessageToUsers(req.OperationID, common.MomentsCommentFriend, pushData, []string{info.UserId}...); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
	}
	http.Success(c, resp)
}

func (s *momentsService) DelMomentsComments(c *gin.Context) {
	var (
		req  model.DelCommentReq
		resp model.MomentsComments
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	if resp, err = momentsRepo.MomentsRepo.MomentsCommentsDel(req, userID); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserPermissions, lang))
		return
	}
	userRepo.UserCache.DelMomentsMessage(userID)
	http.Success(c, resp)
	return
}

func (s *momentsService) LikeMoments(c *gin.Context) {
	var (
		req model.MomentsCommentsLikeReq
		err error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	info, err := momentsRepo.MomentsRepo.MomentsCommentsLikeGet(req.MomentsID, userID)
	pushData := model.MomentsInviteFriendPush{LikeStatus: 2}
	if err == nil {
		pushData.LikeStatus, _ = momentsRepo.MomentsRepo.MomentsCommentsLikeUpdate(info.Status, userID, req.MomentsID)
	} else {
		info, _ = momentsRepo.MomentsRepo.MomentsCommentsLikeAdd(req, userID)
	}
	publisherUser, err := userusecase.UserUseCase.GetBaseInfo(userID)
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrUserNotFound, lang))
		return
	}
	momentsInfo := new(model.MomentsMessage)
	if err = momentsRepo.MomentsRepo.MomentsInfo(req.MomentsID, momentsInfo); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	userRepo.UserCache.DelMomentsMessage(momentsInfo.UserId)
	conversationID := apiChatUseCase.ConversationUseCase.GetConversationID(chatModel.ConversationTypeSingle, userID, momentsInfo.UserId)
	pushData.ConversationID = conversationID
	pushData.Timestamp = time.Now().Unix()
	pushData.PublisherUserID = userID
	pushData.FriendUserID = momentsInfo.UserId
	pushData.MomentsID = req.MomentsID
	pushData.FaceURL = publisherUser.FaceURL
	pushData.PublisherNickname = publisherUser.NickName
	if err = mqtt.SendMessageToUsers(req.OperationID, common.MomentsLikeFriend, pushData, []string{momentsInfo.UserId}...); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
	}
	http.Success(c, nil)
	return
}

func (s *momentsService) MomentsCommentsList(c *gin.Context) {
	var (
		req  model.MomentsCommentsReq
		resp model.MomentsCommentsList
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

	var count int64
	momentsCommentsList, err := momentsRepo.MomentsRepo.MomentsCommentsList(req.MomentsID, &count, resp.Pagination.Offset, resp.Pagination.Limit)
	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	momentsComments := make([]model.CommentsInfo, 0)
	if len(momentsCommentsList) > 0 {
		for _, vv := range momentsCommentsList {
			momentsCommentsInfo := model.CommentsInfo{}
			user, err := userusecase.UserUseCase.GetBaseInfo(vv.UserID)
			if err != nil {
				logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
				http.Failed(c, response.GetError(response.ErrDB, lang))
				return
			}
			momentsCommentsInfo.User.UserId = user.UserId
			momentsCommentsInfo.User.FaceURL = user.FaceURL
			momentsCommentsInfo.User.NickName = user.NickName
			momentsCommentsInfo.MomentsID = vv.MomentsID
			momentsCommentsInfo.ID = vv.ID
			momentsCommentsInfo.IsOwnComment = vv.IsOwnComment
			if vv.ReplyToId != "" {
				momentsCommentsInfo.ReplyToUser.UserId = vv.ReplyUser.UserID
				momentsCommentsInfo.ReplyToUser.FaceURL = vv.ReplyUser.FaceURL
				momentsCommentsInfo.ReplyToUser.NickName = vv.ReplyUser.NickName
			}
			momentsCommentsInfo.Content = vv.Content
			momentsCommentsInfo.Images = strings.Split(vv.Image, ";")
			momentsComments = append(momentsComments, momentsCommentsInfo)
		}
	}
	resp.List = momentsComments
	resp.Count = count
	http.Success(c, resp)
	return

}

func (s *momentsService) MomentsDetail(c *gin.Context) {
	var (
		req     model.MomentsReq
		resp    model.MomentsDetailResp
		moments model.MomentsMessage
		err     error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	resp.Pagination.Check()
	moments, err = momentsRepo.MomentsRepo.GetMomentsDetail(req.MomentsID)
	if resp.Data, err = usecase.MomentsUseCase.MomentsDetail(moments, lang, resp.Offset, resp.Limit); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, err)
		return
	}
	userID := c.GetString("user_id")
	if userID != moments.UserId {
		if moments.CanSee == model.CanSeePrivate {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "朋友圈为私有，您无法查看")
			http.Failed(c, response.GetError(response.ErrBadRequest, lang))
			return
		}
		if moments.CanSee == model.CanSeeFriend && !friendusecase.FriendUseCase.CheckFriend(userID, moments.UserId) {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "朋友圈只能朋友才能查看，您无法查看")
			http.Failed(c, response.GetError(response.ErrBadRequest, lang))
			return
		}
	}
	http.Success(c, resp)
	return
}
