package repo

import (
	"gorm.io/gorm"
	"im/internal/api/discover/model"
	"im/pkg/db"
)

var NewsRepo = new(newRepo)

type newRepo struct{}

func (r *newRepo) List(req model.NewsReq) (news []model.News, count int64, err error) {
	tx := db.DB.Model(&model.News{}).Where("status = ? ", model.StatusOn)
	if req.Title != "" {
		tx = tx.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.OrderBy != "" {
		tx = tx.Order(req.OrderBy)
	}
	err = tx.Order("updated_at DESC").Offset(req.Offset).Limit(req.Limit).Find(&news).Offset(-1).Limit(-1).Count(&count).Error
	return
}

func (r *newRepo) Detail(ID int64) (news model.News, err error) {
	err = db.DB.Model(&model.News{}).Where("status = ? ", model.StatusOn).Where("id = ?", ID).Find(&news).Error
	return
}

func (r *newRepo) ViewTotalInc(ID int64) (err error) {
	err = db.DB.Model(&model.News{}).Where("id = ?", ID).Update("view_total", gorm.Expr("view_total+ ?", 1)).Error
	return
}
