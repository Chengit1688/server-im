package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im/internal/api/discover/model"
	"im/internal/api/discover/repo"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
)

var NewsService = new(newsService)

type newsService struct {
}

func (s *newsService) List(c *gin.Context) {
	var (
		req  model.NewsReq
		resp model.NewsResp
		news []model.News
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
	if news, resp.Count, err = repo.NewsRepo.List(req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("search error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	for _, n := range news {
		l := model.NewsDetailInfo{
			ID:           n.ID,
			CreateUserID: n.CreateUserID,
			Title:        n.Title,
			Content:      n.Content,
			ViewTotal:    n.ViewTotal,
			CategoryID:   n.CategoryID,
			Image:        strings.Split(n.Image, ","),
			Video:        n.Video,
			CreatedAt:    n.CreatedAt,
			UpdatedAt:    n.UpdatedAt,
		}
		if n.Video != "" {
			l.Image = []string{}
		}
		if len(l.Image) > 0 {
			n.Video = ""
		}
		resp.List = append(resp.List, l)
	}

	http.Success(c, resp)
}

func (s *newsService) Detail(c *gin.Context) {
	var (
		req  model.NewsDetailReq
		resp model.NewsDetailInfo
		news model.News
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if news, err = repo.NewsRepo.Detail(req.ID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	if err = repo.NewsRepo.ViewTotalInc(req.ID); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}
	resp = model.NewsDetailInfo{
		ID:           news.ID,
		CreateUserID: news.CreateUserID,
		Title:        news.Title,
		Content:      news.Content,
		ViewTotal:    news.ViewTotal,
		CategoryID:   news.CategoryID,
		Image:        strings.Split(news.Image, ","),
		Video:        news.Video,
		CreatedAt:    news.CreatedAt,
		UpdatedAt:    news.UpdatedAt,
	}

	http.Success(c, resp)
}
