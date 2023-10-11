package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	friendModel "im/internal/api/friend/model"
	"im/internal/api/friend/repo"
	friendusecase "im/internal/api/friend/usecase"
	"im/internal/api/moments/model"
	momentsRepo "im/internal/api/moments/repo"
	"im/pkg/code"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var ContactsService = new(contactsService)

type contactsService struct{}

func (s *contactsService) AddTag(c *gin.Context) {
	var (
		req model.TagReq
		err error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	for _, u := range req.UserID {
		if false == friendusecase.FriendUseCase.CheckFriend(userID, u) {
			http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
			return
		}
	}
	if err = momentsRepo.TagRepo.AddTag(userID, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		if err == code.ErrTagExists {
			http.Failed(c, response.GetError(response.ErrTagExists, lang))
			return
		}
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	for _, u := range req.UserID {

		err = repo.FriendRepo.ChangeFriendVersion(userID, u)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)

		err = friendusecase.FriendUseCase.UpdateFriend(userID, u)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
	}

	http.Success(c)
}

func (s *contactsService) AddFriendTag(c *gin.Context) {
	var (
		req model.TagAddFriendReq
		err error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	hadFriend := friendModel.Friend{}
	if err = db.Info(&hadFriend, friendModel.Friend{
		OwnerUserID:  userID,
		FriendUserID: req.UserID,
		Status:       1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if err = momentsRepo.TagRepo.AddFriendTag(userID, req.UserID, req.TagID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		if err == code.ErrTagFriendExists {
			http.Failed(c, response.GetError(response.ErrTagFriendExists, lang))
			return
		}
		http.Failed(c, response.GetError(response.ErrTagNotExists, lang))
		return
	}

	err = repo.FriendRepo.ChangeFriendVersion(userID, req.UserID)
	logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)

	err = friendusecase.FriendUseCase.UpdateFriend(userID, req.UserID)
	logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)

	http.Success(c)
}

func (s *contactsService) CheckFriendTag(c *gin.Context) {
	var (
		req model.TagAddFriendReq
		err error
		b   bool
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	hadFriend := friendModel.Friend{}
	if err = db.Info(&hadFriend, friendModel.Friend{
		OwnerUserID:  userID,
		FriendUserID: req.UserID,
		Status:       1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if b, err = momentsRepo.TagRepo.CheckFriendTag(userID, req.UserID, req.TagID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrTagNotExists, lang))
		return
	}
	if !b {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrTagFriendNotExists, lang))
		return
	}

	http.Success(c)
}

func (s *contactsService) FetchFriendTag(c *gin.Context) {
	var (
		req  model.FetchFriendTag
		resp model.TagAddFriendResp
		err  error
		tags []model.ContactsTag
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	hadFriend := friendModel.Friend{}
	if err = db.Info(&hadFriend, friendModel.Friend{
		OwnerUserID:  userID,
		FriendUserID: req.UserID,
		Status:       1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrFriendNotExist, lang))
		return
	}
	if tags, err = momentsRepo.TagRepo.GetFriendTag(userID, req.UserID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrTagNotExists, lang))
		return
	}
	for _, tag := range tags {
		resp.List = append(resp.List, model.TagListInfo{
			Title: tag.Title,
			TagID: tag.ID,
		})
	}

	http.Success(c, resp)
}

func (s *contactsService) UpdateTag(c *gin.Context) {
	var (
		req model.TagReq
		err error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.TagID == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "tagID不能为空")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	if _, err = momentsRepo.TagRepo.EditeTag(req.TagID, userID, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	for _, u := range req.UserID {

		err = repo.FriendRepo.ChangeFriendVersion(userID, u)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)

		err = friendusecase.FriendUseCase.UpdateFriend(userID, u)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
	}

	http.Success(c)
}

func (s *contactsService) DeleteTag(c *gin.Context) {
	var (
		req   model.CommonIDReq
		data  model.ContactsTag
		err   error
		count int64
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	if count, err = momentsRepo.TagRepo.TagCountByCreatorID(userID, req.ID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	if count > 0 {
		http.Failed(c, response.GetError(response.ErrTagFriendExists, lang))
		return
	}
	if err = momentsRepo.TagRepo.DeleteTagByID(req.ID, userID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	if data, err = momentsRepo.TagRepo.GetTag(req.ID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	userIDs := strings.Split(data.FriendUserID, ",")
	for _, u := range userIDs {

		err = repo.FriendRepo.ChangeFriendVersion(userID, u)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)

		err = friendusecase.FriendUseCase.UpdateFriend(userID, u)
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
	}

	http.Success(c)
}

func (s *contactsService) ListTag(c *gin.Context) {
	var (
		req  model.TagListReq
		resp model.TagListResp
		tags []model.ContactsTag
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	req.Pagination.Check()
	if tags, resp.Count, err = momentsRepo.TagRepo.ListTagByCreatorID(userID, req); err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrDB, lang))
			return
		}
	}
	for _, tag := range tags {
		i := model.TagListInfo{
			TagID: tag.ID,
			Title: tag.Title,
		}
		friendUserID := strings.Split(tag.FriendUserID, ",")
		i.UserTotal = len(friendUserID)
		for _, f := range friendUserID {
			if f == "0" || f == "" {
				i.UserTotal -= 1
				continue
			}
		}

		resp.List = append(resp.List, i)
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *contactsService) TagDetail(c *gin.Context) {
	var (
		req  model.CommonIDReq
		resp model.TagDetailResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	userID := c.GetString("user_id")
	if resp, err = momentsRepo.TagRepo.TagDetailByID(userID, req.ID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	http.Success(c, resp)
}
