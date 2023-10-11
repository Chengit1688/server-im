package repo

import (
	apiNews "im/internal/api/discover/model"
	"im/internal/cms_api/discover/model"
	"im/pkg/db"
	"strings"
	"time"
)

var NewRepo = new(newRepo)

type newRepo struct{}

func (r *newRepo) AddNews(createUserID string, req model.AddNewsReq) (data apiNews.News, err error) {
	data = apiNews.News{
		CreateUserID: createUserID,
		Title:        req.Title,
		Content:      req.Content,
		CategoryID:   req.CategoryID,
		Status:       req.Status,
		Image:        strings.Join(req.Image, ","),
		Video:        req.Video,
		CommonModel: db.CommonModel{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
	err = db.DB.Model(&apiNews.News{}).Create(&data).Error
	return data, err
}
func (r *newRepo) UpdateNews(ID int64, req model.AddNewsReq) (data apiNews.News, err error) {
	data = apiNews.News{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		Image:      strings.Join(req.Image, ","),
		Video:      req.Video,
		CommonModel: db.CommonModel{
			UpdatedAt: time.Now().Unix(),
		},
	}
	err = db.DB.Model(&apiNews.News{}).Where("id = ?", ID).Updates(&data).Error
	return data, err
}

func (r *newRepo) DeleteNews(deleteUserID string, id int64) (data apiNews.News, err error) {
	data = apiNews.News{
		Status:       apiNews.StatusOff,
		DeleteUserID: deleteUserID,
		CommonModel: db.CommonModel{
			UpdatedAt: time.Now().Unix(),
		},
	}
	err = db.DB.Model(&apiNews.News{}).Where("id = ?", id).Updates(&data).Error
	return data, err
}

func (r *newRepo) List(req model.ListNewsReq) (news []apiNews.News, count int64, err error) {
	tx := db.DB.Model(&apiNews.News{}).Preload("CreatorUser").Where("status = ?", apiNews.StatusOn)
	if req.Title != "" {
		tx = tx.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.CreateUserID != "" {
		tx = tx.Where("create_user_id = ?", req.CreateUserID)
	}
	if req.OrderBy != "" {
		tx = tx.Order(req.OrderBy)
	}
	err = tx.Order("updated_at DESC").Offset(req.Offset).Limit(req.Limit).Find(&news).Offset(-1).Limit(-1).Count(&count).Error
	return
}
