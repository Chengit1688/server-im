package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apiNewsModel "im/internal/api/discover/model"
	"im/internal/cms_api/discover/model"
	"im/internal/cms_api/discover/repo"
	"im/pkg/code"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var NewsService = new(newService)

type newService struct{}

func (s *newService) Add(c *gin.Context) {
	var (
		err error
		req model.AddNewsReq
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	userID := c.GetString("user_id")
	if req.Status == 0 {
		req.Status = apiNewsModel.StatusOn
	}
	if _, err = repo.NewRepo.AddNews(userID, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	http.Success(c)
}

func (s *newService) Delete(c *gin.Context) {
	var (
		err error
		req model.DeleteNewsReq
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	userID := c.GetString("user_id")
	if _, err = repo.NewRepo.DeleteNews(userID, req.ID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	http.Success(c)
}

func (s *newService) Update(c *gin.Context) {
	var (
		err error
		req model.AddNewsReq
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.ID == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Status == 0 {
		req.Status = apiNewsModel.StatusOn
	}
	if _, err = repo.NewRepo.UpdateNews(req.ID, req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	http.Success(c)
}

func (s *newService) List(c *gin.Context) {
	var (
		req  model.ListNewsReq
		resp model.ListNewsResp
		news []apiNewsModel.News
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
	if news, resp.Count, err = repo.NewRepo.List(req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("search error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	for _, n := range news {
		resp.List = append(resp.List, model.NewsInfo{
			ID:             n.ID,
			CreateUserID:   n.CreateUserID,
			Title:          n.Title,
			Content:        n.Content,
			ViewTotal:      n.ViewTotal,
			CategoryID:     n.CategoryID,
			Image:          strings.Split(n.Image, ","),
			Video:          n.Video,
			CreateNickname: n.CreatorUser.NickName,
			CreateAccount:  n.CreatorUser.Account,
			CreatedAt:      n.CreatedAt,
			UpdatedAt:      n.UpdatedAt,
		})
	}
	http.Success(c, resp)
}
